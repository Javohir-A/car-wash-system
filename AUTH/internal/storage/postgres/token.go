package postgres

import (
	"auth-service/config"
	"auth-service/genprotos/auth"
	"auth-service/internal/storage"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type TokenSQL struct {
	userManager storage.UserManagement
	db          *sql.DB
	sql         sq.StatementBuilderType
	redis       *redis.Client
}

func NewTokenSQL(db *sql.DB, redis *redis.Client, userManager storage.UserManagement) *TokenSQL {
	return &TokenSQL{
		db:          db,
		sql:         sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
		redis:       redis,
		userManager: userManager,
	}
}

type Claims struct {
	User *auth.UserResponse `json:"user"`
	jwt.RegisteredClaims
}

func (t *TokenSQL) GenerateToken(ctx context.Context, userId string) (string, error) {
	user, err := t.userManager.GetUserByID(ctx, userId)
	if err != nil {
		return "", err
	}
	log.Println(user)
	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		return "", err
	}

	jwtKey := cnf.JWT.SecretKey

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})

	accessStr, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	refreshStr, err := refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	tx, err := t.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := t.sql.
		Insert("tokens").
		Columns("user_id", "access_token", "refresh_token").
		Values(userId, accessStr, refreshStr)

	sqlString, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	_, err = tx.ExecContext(ctx, sqlString, args...)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	cacheAKey := fmt.Sprintf("accessToken:%s", userId)
	cacheRKey := fmt.Sprintf("refreshToken:%s", userId)

	cachedAData, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}

	cachedRData, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}

	if err := t.redis.Set(ctx, cacheAKey, cachedAData, time.Hour*24).Err(); err != nil {
		return "", err
	}

	if err := t.redis.Set(ctx, cacheRKey, cachedRData, time.Hour*24*7).Err(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", accessStr, refreshStr), nil
}

func (t *TokenSQL) InvalidateToken(ctx context.Context, tokenStr string) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func ExtractClaims(tokenStr string) (*Claims, error) {
	var claims *Claims
	cnf := config.NewConfig()
	cnf.Load()

	jwtKey := cnf.JWT.SecretKey

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
