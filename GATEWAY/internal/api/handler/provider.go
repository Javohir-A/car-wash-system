package handler

import (
	"context"
	pb "gateway/genproto/providers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProviderManagementHandler struct {
	pro pb.ProvidersClient
}

func NewProviderManagementHandler(client pb.ProvidersClient) *ProviderManagementHandler {
	return &ProviderManagementHandler{
		pro: client,
	}
}

// RegisterProviderHandler godoc
// @Summary Register a new provider
// @Description Register a new provider with the given details
// @Tags Provider
// @Accept json
// @Produce json
// @Param request body providers.NewProvider true "RegisterProvider Request"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/provider/register [post]
func (h *ProviderManagementHandler) RegisterProvider(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.NewProvider
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := h.pro.CreateProvider(ctx, &req) // Assuming the gRPC method is still named `CreateProvider`
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't register provider"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// SearchProvidersHandler godoc
// @Summary Search for providers
// @Description Search for providers based on filters such as name, average rating, and creation date
// @Tags Provider
// @Accept json
// @Produce json
// @Param request body providers.Filter true "SearchProviders Request"
// @Security BearerAuth
// @Success 200 {object} providers.SearchResp
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/provider/search [post]
func (h *ProviderManagementHandler) SearchProviders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.Filter
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.pro.SearchProviders(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't search providers"})
		return
	}
	c.JSON(http.StatusOK, res)
}
