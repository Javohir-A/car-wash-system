package handler

import (
	"context"
	"gateway/genproto/auth"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

type UserManagementHandler struct {
	user auth.UserManagementServiceClient
}

func NewUserManagementHandler(client auth.UserManagementServiceClient) *UserManagementHandler {
	return &UserManagementHandler{
		user: client,
	}
}

// UserCreateHandler godoc
// @Summary UserCreate a new user
// @Description Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body auth.UserRequest true "CreateUser Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/ [post]
func (u *UserManagementHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req auth.UserRequest
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := u.user.CreateUser(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't create user"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// GetUserById godoc
// @Summary Get a user by their id
// @Description Get user by id
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Get user by id Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/{id} [get]
func (u *UserManagementHandler) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := c.Param("id")
	log.Println(id)

	if err := uuid.Validate(id); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be uuid"})
		return
	}

	res, err := u.user.GetUserById(ctx, &auth.UserIdRequest{Id: id})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch user data"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// DeleteUserByID godoc
// @Summary Delete a user by their id
// @Description Delete user by id
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Delete user by id Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/{id} [delete]
func (u *UserManagementHandler) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")

	if err := uuid.Validate(id); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be uuid"})
		return
	}

	_, err := u.user.DeleteUser(ctx, &auth.UserIdRequest{Id: id})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch user data"})
		return
	}
	c.JSON(http.StatusOK, "user deleted successfully")
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user
// @Tags User
// @Accept json
// @Produce json
// @Param request body auth.UserRequest true "UpdateUser Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/ [put]
func (u *UserManagementHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var updatedUser auth.UserRequest
	if err := c.BindJSON(&updatedUser); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := uuid.Validate(updatedUser.Id); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be uuid"})
		return
	}

	res, err := u.user.UpdateUser(ctx, &updatedUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't update user data"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetUsers godoc
// @Summary get users by page and limit
// @Description Get Users
// @Tags User
// @Accept json
// @Produce json
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/ [get]
func (u *UserManagementHandler) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	l := c.Query("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	p := c.Query("page")
	page, err := strconv.Atoi(p)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}

	res, err := u.user.GetUsers(ctx, &auth.GetUsersRequest{PageNumber: int32(page), PageSize: int32(limit)})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, res)
}
