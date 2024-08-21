package postgres

import (
	"auth-service/genprotos/auth"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserManagementImpl struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
	redis      *redis.Client
}

func NewUserManagementSQL(db *sql.DB, redis *redis.Client) *UserManagementImpl {
	return &UserManagementImpl{
		db:         db,
		sqlBuilder: sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
		redis:      redis,
	}
}

func (um *UserManagementImpl) CreateUser(ctx context.Context, user *auth.UserRequest) (*auth.UserResponse, error) {
	tx, err := um.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user.Id = uuid.NewString()

	query := um.sqlBuilder.
		Insert("users").
		Columns("id", "first_name", "last_name", "phone_number", "username", "password_hash", "email", "role").
		Values(user.Id, user.FirstName, user.LastName, user.PhoneNumber, user.Username, user.Password, user.Email, user.Role).
		Suffix("RETURNING created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, sql, args...)
	res := &auth.UserResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
	}
	var createdAt, updatedAt string
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		log.Println("failed to scan user")
		return nil, err
	}
	createdTimeType, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, nil
	}
	updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, nil
	}
	res.CreatedAt = timestamppb.New(createdTimeType)
	res.UpdatedAt = timestamppb.New(updatedTimeType)

	cacheKey := fmt.Sprintf("user:%v", user.Id)

	value, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	_, err = um.redis.Set(ctx, cacheKey, value, time.Hour*2).Result()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}
func (um *UserManagementImpl) GetUserByID(ctx context.Context, id string) (*auth.UserResponse, error) {
	cacheKey := fmt.Sprintf("user:%v", id)
	value, err := um.redis.Get(ctx, cacheKey).Result()

	var user auth.UserResponse
	fmt.Printf("Cache value: %s\n", value)
	if err == redis.Nil || value == "" {
		query := um.sqlBuilder.
			Select("id", "first_name", "last_name", "phone_number", "username", "email", "password_hash", "role", "created_at", "updated_at").
			From("users").
			Where(sq.And{
				sq.Eq{"id": id},
				sq.Eq{"deleted_at": nil},
			})

		sqlString, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}

		row := um.db.QueryRowContext(ctx, sqlString, args...)
		var createdAt, updatedAt string
		err = row.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&createdAt,
			&updatedAt,
		)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		} else if err != nil {
			return nil, err
		}

		createdTimeType, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = timestamppb.New(createdTimeType)
		user.UpdatedAt = timestamppb.New(updatedTimeType)

		cacheData, err := json.Marshal(&user)
		if err != nil {
			log.Println("caching data failed")
			return nil, err
		}

		_, err = um.redis.Set(ctx, cacheKey, cacheData, time.Hour*2).Result()
		if err != nil {
			log.Println("storing data failed")
			return nil, err
		}

		return &user, nil
	}

	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
	}

	return &user, nil
}
func (um *UserManagementImpl) UpdateUser(ctx context.Context, user *auth.UserRequest) (*auth.UserResponse, error) {
	if user == nil {
		return nil, errors.New("user request cannot be nil")
	}

	tx, err := um.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := um.sqlBuilder.Update("users").
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("phone_number", user.PhoneNumber).
		Set("username", user.Username).
		Set("role", user.Role).
		Set("password_hash", user.Password).
		Where(sq.And{
			sq.Eq{"id": user.Id},
			sq.Eq{"deleted_at": nil},
		}).Suffix("RETURNING created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	log.Printf("Executing SQL: %s", sql)
	row := tx.QueryRowContext(ctx, sql, args...)

	res := &auth.UserResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
	}

	var createdAt, updatedAt string
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		log.Println("Failed to scan user")
		return nil, err
	}

	createdTimeType, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}
	res.CreatedAt = timestamppb.New(createdTimeType)
	res.UpdatedAt = timestamppb.New(updatedTimeType)

	cacheKey := fmt.Sprintf("user:%v", user.Id)

	cacheData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	_, err = um.redis.Set(ctx, cacheKey, cacheData, time.Hour*24).Result()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}

func (um *UserManagementImpl) DeleteUser(ctx context.Context, id string) error {
	tx, err := um.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := um.sqlBuilder.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Eq{"deleted_at": nil},
		},
		)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, sql, args...)
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%v", id)

	um.redis.Del(ctx, cacheKey)

	tx.Commit()
	return nil
}
func (um *UserManagementImpl) ListUsers(ctx context.Context, pageNumber, pageSize int32) ([]*auth.UserResponse, error) {
	fmt.Println(pageNumber, pageSize)
	if pageNumber < 1 {
		return nil, fmt.Errorf("invalid page number: %d", pageNumber)
	}
	if pageSize < 1 {
		return nil, fmt.Errorf("invalid page size: %d", pageSize)
	}

	var users []*auth.UserResponse
	offset := (pageNumber - 1) * pageSize

	query := um.sqlBuilder.
		Select("id", "first_name", "last_name", "phone_number", "username", "email", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"deleted_at": nil}).
		Limit(uint64(pageSize)).
		Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := um.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user auth.UserResponse
		var createdAt, updatedAt string
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Username, &user.Email, &user.Role, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		createdTimeType, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, nil
		}
		updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, nil
		}
		user.CreatedAt = timestamppb.New(createdTimeType)
		user.UpdatedAt = timestamppb.New(updatedTimeType)

		users = append(users, &user)
	}

	return users, nil
}

func (um *UserManagementImpl) SearchUsers(ctx context.Context, queryStr string) ([]*auth.UserResponse, error) {
	var users []*auth.UserResponse

	query := um.sqlBuilder.
		Select("*").
		From("users").
		Where(sq.Or{
			sq.Like{"first_name": "%" + queryStr + "%"},
			sq.Like{"last_name": "%" + queryStr + "%"},
			sq.Like{"username": "%" + queryStr + "%"},
			sq.Like{"email": "%" + queryStr + "%"},
		}).
		Limit(10)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql, args)
	rows, err := um.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumer string
	for rows.Next() {
		var user auth.UserResponse
		var createdAt, updatedAt string

		createdTimeType, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, nil
		}
		updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, nil
		}
		user.CreatedAt = timestamppb.New(createdTimeType)
		user.UpdatedAt = timestamppb.New(updatedTimeType)
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&consumer,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (um *UserManagementImpl) GetUserByUsernameOrEmail(ctx context.Context, email, username string) (*auth.UserResponse, error) {
	builder := um.sqlBuilder.Select(
		"id",
		"first_name",
		"last_name",
		"phone_number",
		"username",
		"email",
		"password_hash",
		"role",
		"created_at",
		"updated_at",
	).From("users")

	filter := sq.Eq{}

	filter["email"] = ""
	filter["username"] = ""

	if email != "" {
		filter["email"] = email
	}
	if username != "" {
		filter["username"] = username
	}

	query, args, err := builder.Where(sq.Or{
		sq.Eq{"username": username},
	}).ToSql()

	if err != nil {
		return nil, err
	}

	var user auth.UserResponse

	row := um.db.QueryRowContext(ctx, query, args...)
	var createdAt, updatedAt string

	err = row.Scan(&user.Id,
		&user.FirstName, &user.LastName,
		&user.PhoneNumber, &user.Username,
		&user.Email, &user.Password,
		&user.Role, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}
	createdTimeType, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, nil
	}
	updatedTimeType, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, nil
	}
	user.CreatedAt = timestamppb.New(createdTimeType)
	user.UpdatedAt = timestamppb.New(updatedTimeType)

	return &user, nil
}
