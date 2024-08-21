package handler

import (
	"context"
	"gateway/genproto/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SudoHandler struct {
	user auth.UserManagementServiceClient
}

var existingRoles = map[string]bool{
	"admin":       true,
	"super admin": true,
	"provider":    true,
	"user":        true,
}

func NewSudoHandler(user auth.UserManagementServiceClient) *SudoHandler {
	return &SudoHandler{
		user: user,
	}
}

// ChangeUserRole godoc
// @Summary Change user role, only for sudo
// @Description Give permmisson to a user
// @Tags Sudo
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Param role query string true "role"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sudo/change-role/{user_id} [post]
func (s *SudoHandler) ChangeUserRole(c *gin.Context) {
	id := c.Param("user_id")
	role := c.Query("role")

	if err := uuid.Validate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id: must be uuid"})
		return
	}
	log.Println(id, role)
	if !existingRoles[role] {
		log.Printf("No such role as %s", role)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such role!"})
		return
	}

	grpcReq := auth.UserRequest{
		Id:   id,
		Role: role,
	}

	grpcRes, err := s.user.UpdateUser(context.TODO(), &grpcReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, grpcRes)
}
