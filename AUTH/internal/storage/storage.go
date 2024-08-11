package storage

import (
	"auth-service/genprotos/auth"
	"context"
)

type Storage interface {
	UserManagement() UserManagement
	// Authentification() Authentication
	TokenManagement() TokenManagement
}

type UserManagement interface {
	CreateUser(ctx context.Context, user *auth.UserRequest) (*auth.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*auth.UserResponse, error)
	UpdateUser(ctx context.Context, user *auth.UserRequest) (*auth.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int32) ([]*auth.UserResponse, error)
	SearchUsers(ctx context.Context, query string) ([]*auth.UserResponse, error)
	GetUserByUsernameOrEmail(ctx context.Context, email, username string) (*auth.UserResponse, error)
}

// type Authentication interface {
// }

type TokenManagement interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	// InvalidateToken(ctx context.Context, token string) error
	// RefreshToken(ctx context.Context, refreshToken string) (string, error)
	// RevokeTokens(ctx context.Context, userID string) error
}
