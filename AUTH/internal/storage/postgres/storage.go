package postgres

import (
	"auth-service/internal/storage"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
)

type StorageSQL struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
	redis      *redis.Client
}

func NewStorageSQL(db *sql.DB, redis *redis.Client) *StorageSQL {
	return &StorageSQL{
		db:         db,
		sqlBuilder: sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
		redis:      redis,
	}
}

func (s *StorageSQL) UserManagement() storage.UserManagement {
	return NewUserManagementSQL(s.db, s.redis)
}

// func (s *StorageSQL) Authentification() Authentication {
// 	return NewAuthentication(s.db)
// }

func (s *StorageSQL) TokenManagement() storage.TokenManagement {
	return NewTokenSQL(s.db, s.redis, s.UserManagement())
}
