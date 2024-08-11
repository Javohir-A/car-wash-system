package repository

import (
	"auth-service/config"
	"auth-service/internal/service"
	"auth-service/internal/storage/postgres"

	"database/sql"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

func SetupAuthService(db *sql.DB, rClient *redis.Client, cnf *config.Config, logger slog.Logger) (*service.AuthServiceImpl, error) {

	storage := postgres.NewStorageSQL(db, rClient)

	return service.NewAuthService(storage, &logger), nil
}
