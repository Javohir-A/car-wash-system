package service

import (
	"auth-service/genprotos/auth"
	"auth-service/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServiceImpl struct {
	storage storage.Storage
	auth.UnimplementedAuthServiceServer
	auth.UnimplementedUserManagementServiceServer
	logger *slog.Logger
}

func NewAuthService(storage storage.Storage,
	logger *slog.Logger,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (a *AuthServiceImpl) InvalidateToken(ctx context.Context, req *auth.InvalidateTokenRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.Println(req)
	res, err := a.storage.UserManagement().GetUserByUsernameOrEmail(ctx, "", req.Username)
	log.Println(req)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}
	if res == nil {
		a.logger.Warn("no user found")
		return nil, fmt.Errorf("no user found")
	}
	hash := res.GetPassword()

	err, ok := ComparePassword(hash, req.Password)
	if err != nil {
		if !ok {
			a.logger.Warn("invalid passord")
			return nil, fmt.Errorf("invalidPassord")
		}
		a.logger.Error(err.Error())
		return nil, err
	}

	token, err := a.storage.TokenManagement().GenerateToken(ctx, res.GetId())
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(token, ":")

	return &auth.LoginResponse{
		AccessToken:            tokens[0],
		RefreshToken:           tokens[1],
		AccessTokenExpiration:  timestamppb.New(time.Now().Add(time.Hour * 24)),
		RefreshTokenExpiration: timestamppb.New(time.Now().Add(time.Hour * 24 * 7)),
	}, nil
}
func (a *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.LoginResponse, error) {
	return nil, nil
}

// TODO: should be checked the username is unique
func (a *AuthServiceImpl) CreateUser(ctx context.Context, req *auth.UserRequest) (*auth.UserResponse, error) {
	_, err := a.storage.UserManagement().GetUserByUsernameOrEmail(ctx, req.GetEmail(), req.GetUsername())
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("username or email exists")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	req.Password = hash

	res, err := a.storage.UserManagement().CreateUser(ctx, req)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func (a *AuthServiceImpl) DeleteUser(ctx context.Context, req *auth.UserIdRequest) (*emptypb.Empty, error) {
	err := a.storage.UserManagement().DeleteUser(ctx, req.Id)
	if err != nil {
		a.logger.Error(err.Error())
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServiceImpl) GetUserById(ctx context.Context, req *auth.UserIdRequest) (*auth.UserResponse, error) {
	res, err := a.storage.UserManagement().GetUserByID(ctx, req.Id)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func (a *AuthServiceImpl) GetUsers(req *auth.GetUsersRequest, stream auth.UserManagementService_GetUsersServer) error {
	pageSize := int32(10)
	pageNumber := req.GetPageNumber()

	for i := int32(1); i < pageNumber; i++ {
		users, err := a.storage.UserManagement().ListUsers(stream.Context(), i, pageSize)
		if err != nil {
			return err
		}

		for _, user := range users {
			if err := stream.Send(user); err != nil {
				return err
			}
		}

		if len(users) < int(pageSize) {
			break
		}
	}
	return nil
}

func (a *AuthServiceImpl) UpdateUser(ctx context.Context, req *auth.UserRequest) (*auth.UserResponse, error) {
	user, err := a.storage.UserManagement().GetUserByID(ctx, req.GetId())
	if err != nil {
		log.Println("failed to get user to compare with updated user data: ", err)
		return nil, err
	}

	if req.Username != "" || req.Email != "" {
		_, err := a.storage.UserManagement().GetUserByUsernameOrEmail(ctx, req.GetEmail(), req.GetUsername())
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("username or email exists")
		}
	}

	updatedFieldsCount := 0

	if req.Password != "" {
		req.Password, err = HashPassword(req.Password)
		if err != nil {
			a.logger.Error(err.Error())
			return nil, err
		}
		updatedFieldsCount++
	}

	if req.Email == "" {
		req.Email = user.Email
		updatedFieldsCount++
	}
	if req.FirstName == "" {
		req.FirstName = user.FirstName
		updatedFieldsCount++
	}
	if req.LastName == "" {
		req.LastName = user.LastName
		updatedFieldsCount++
	}
	if req.PhoneNumber == "" {
		req.PhoneNumber = user.PhoneNumber
		updatedFieldsCount++
	}
	if req.Role == "" {
		req.Role = user.Role
		updatedFieldsCount++
	}
	if req.Username == "" {
		req.Username = user.Username
		updatedFieldsCount++
	}
	if req.Password == "" {
		req.Password = user.Password
	}

	if updatedFieldsCount == 0 {
		return nil, fmt.Errorf("no fields provided to update")
	}

	res, err := a.storage.UserManagement().UpdateUser(ctx, req)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func ComparePassword(hashedPassword, password string) (error, bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err, false
	}
	return nil, true
}

func HashPassword(passord string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passord), 10)
	if err != nil {
		return "", err
	}
	return string(hash), err
}
