package handler

import (
	"context"
	pb "gateway/genproto/services" // Import the generated gRPC service proto
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ServiceManagementHandler struct {
	serviceClient pb.ServicesClient
}

func NewServiceManagementHandler(client pb.ServicesClient) *ServiceManagementHandler {
	return &ServiceManagementHandler{
		serviceClient: client,
	}
}

// CreateServiceHandler godoc
// @Summary Create a new service
// @Description Create a new service in the system
// @Tags Service
// @Accept json
// @Produce json
// @Param request body services.NewService true "CreateService Request"
// @Success 201 {object} services.CreateResp
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /services [post]
func (h *ServiceManagementHandler) CreateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.NewService
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.serviceClient.CreateService(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create service"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListServicesHandler godoc
// @Summary List all services
// @Description Retrieve a list of all services with pagination
// @Tags Service
// @Accept json
// @Produce json
// @Param request body services.Pagination true "Pagination Request"
// @Success 200 {object} services.ServicesList
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /services [get]
func (h *ServiceManagementHandler) ListServices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.Pagination
	if err := c.BindQuery(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	res, err := h.serviceClient.ListServices(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't list services"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetServiceByIDHandler godoc
// @Summary Get a specific service by ID
// @Description Retrieve a service by its ID
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} services.Service
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /services/{id} [get]
func (h *ServiceManagementHandler) GetServiceByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	serviceID := c.Param("id")
	if serviceID == "" {
		log.Println("Missing service ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing service ID"})
		return
	}

	req := &pb.ID{Id: serviceID}

	res, err := h.serviceClient.GetService(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get service"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateServiceHandler godoc
// @Summary Update a specific service by ID
// @Description Update the details of a service by its ID
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param request body services.NewData true "UpdateService Request"
// @Success 200 {object} services.UpdateResp
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /services/{id} [put]
func (h *ServiceManagementHandler) UpdateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	serviceID := c.Param("id")
	if serviceID == "" {
		log.Println("Missing service ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing service ID"})
		return
	}

	var req pb.NewData
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Id = serviceID

	res, err := h.serviceClient.UpdateService(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update service"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteServiceHandler godoc
// @Summary Delete a specific service by ID
// @Description Delete a service from the system by its ID
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /services/{id} [delete]
func (h *ServiceManagementHandler) DeleteService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	serviceID := c.Param("id")
	if serviceID == "" {
		log.Println("Missing service ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing service ID"})
		return
	}

	req := &pb.ID{Id: serviceID}
	_, err := h.serviceClient.DeleteService(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete service"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
