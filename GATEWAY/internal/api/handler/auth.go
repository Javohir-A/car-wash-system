package handler

import (
	"context"
	"gateway/genproto/auth"
	"gateway/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang/protobuf/ptypes/timestamp"
)

type AuthHandler struct {
	auth auth.AuthServiceClient
	user auth.UserManagementServiceClient
}

func NewAuthHandler(authClient auth.AuthServiceClient, userClient auth.UserManagementServiceClient) *AuthHandler {
	return &AuthHandler{
		auth: authClient,
		user: userClient,
	}
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with an email and password
// @Tags User Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (a *AuthHandler) Register(c *gin.Context) {
	ctx := context.Background()
	var httpReq models.RegisterRequest

	if err := c.BindJSON(&httpReq); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Println(httpReq)

	grpcReq := auth.UserRequest{
		FirstName:   httpReq.FirstName,
		LastName:    httpReq.LastName,
		PhoneNumber: httpReq.PhoneNumber,
		Username:    httpReq.Username,
		Password:    httpReq.Password,
		Email:       httpReq.Email,
		Role:        "user",
	}

	log.Println(httpReq)
	// log.Println(&grpcReq)
	res, err := a.user.CreateUser(ctx, &grpcReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// LoginHandler godoc
// @Summary User login
// @Description Log in a user with email and password
// @Tags User Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (a *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req auth.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "login failed"})
		return
	}

	res, err := a.auth.Login(ctx, &req)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) LogOut(c *gin.Context) {

}
