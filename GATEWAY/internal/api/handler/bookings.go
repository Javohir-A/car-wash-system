package handler

import (
	"context"
	pb "gateway/genproto/bookings"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingClient pb.BookingsClient
}

func NewBookingHandler(client pb.BookingsClient) *BookingHandler {
	return &BookingHandler{
		bookingClient: client,
	}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking with the given details
// @Tags Bookings
// @Accept json
// @Produce json
// @Param request body bookings.NewBooking true "CreateBooking Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /booking [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.NewBooking
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.bookingClient.CreateBooking(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create booking"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListBookings godoc
// @Summary List all bookings
// @Description Retrieve a list of bookings with optional filtering
// @Tags Bookings
// @Accept json
// @Produce json
// @Param limit query string true "limit"
// @Param page query string true "page"
// @Success 200 {object} bookings.BookingsList
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /booking/search [post]
func (h *BookingHandler) ListBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	l := c.Query("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
	}
	p := c.Query("page")
	page, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
	}

	req := pb.Pagination{
		Limit: int32(limit),
		Page:  int32(page),
	}

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.bookingClient.ListBookings(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't list bookings"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetBooking godoc
// @Summary Get booking by ID
// @Description Retrieve booking details by its ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} bookings.Booking
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /booking/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")
	req := &pb.ID{Id: id}

	res, err := h.bookingClient.GetBooking(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't retrieve booking"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateBooking godoc
// @Summary Update booking by ID
// @Description Update the details of a booking by its ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param request body bookings.NewData true "UpdateBooking Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /booking/{id} [put]
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")
	var req pb.NewData
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.Id = id

	res, err := h.bookingClient.UpdateBooking(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update booking"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteBooking godoc
// @Summary Delete booking by ID
// @Description Delete a booking by its ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /booking/{id} [delete]
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")
	req := &pb.ID{Id: id}

	_, err := h.bookingClient.CancelBooking(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete booking"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
