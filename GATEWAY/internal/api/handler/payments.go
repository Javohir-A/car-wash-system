package handler

import (
	"context"
	pb "gateway/genproto/payments"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentClient pb.PaymentsClient
}

func NewPaymentHandler(client pb.PaymentsClient) *PaymentHandler {
	return &PaymentHandler{
		paymentClient: client,
	}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a new payment with the given details
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body payments.NewPayment true "CreatePayment Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.NewPayment
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.paymentClient.CreatePayment(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create payment"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListPayments godoc
// @Summary List all payments
// @Description Retrieve a list of payments with optional filtering
// @Tags Payments
// @Accept json
// @Produce json
// @Param limit query string true "limit"
// @Param page query string true "page"
// @Success 200 {object} payments.PaymentsList
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments [get]
func (h *PaymentHandler) ListPayments(c *gin.Context) {
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
		Page:  int32(page),
		Limit: int32(limit),
	}
	res, err := h.paymentClient.ListPayments(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't list payments"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetPayment godoc
// @Summary Get payment by ID
// @Description Retrieve payment details by its ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} payments.Payment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")
	req := &pb.ID{Id: id}

	res, err := h.paymentClient.GetPayment(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't retrieve payment"})
		return
	}
	c.JSON(http.StatusOK, res)
}
