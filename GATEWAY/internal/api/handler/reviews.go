package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "gateway/genproto/reviews"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewClient pb.ReviewsClient
}

func NewReviewHandler(client pb.ReviewsClient) *ReviewHandler {
	return &ReviewHandler{
		reviewClient: client,
	}
}

// CreateReview godoc
// @Summary Create a new review
// @Description Create a new review with the given details
// @Tags Reviews
// @Accept json
// @Produce json
// @Param request body reviews.NewReview true "CreateReview Request"
// @Security BearerAuth
// @Success 201 {object} reviews.CreateResp
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req pb.NewReview
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.reviewClient.CreateReview(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create review"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListReviews godoc
// @Summary List all reviews
// @Description Retrieve a list of reviews with optional filtering
// @Tags Reviews
// @Accept json
// @Produce json
// @Param limit query string true "limit"
// @Param page query string true "page"
// @Security BearerAuth
// @Success 200 {object} reviews.ReviewsList
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	l := c.Query("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	p := c.Query("page")
	page, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	req := pb.Pagination{
		Limit: int32(limit),
		Page:  int32(page),
	}

	res, err := h.reviewClient.ListReviews(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't list reviews"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetReview godoc
// @Summary Get review by ID
// @Description Retrieve review details by its ID
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} reviews.Review
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews/{id} [get]
// func (h *ReviewHandler) GetReview(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	id := c.Param("id")
// 	req := &pb.ID{Id: id}

// 	res, err := h.reviewClient.GetReview(ctx, req)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't retrieve review"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// UpdateReview godoc
// @Summary Update review by ID
// @Description Update the details of a review by its ID
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Param request body reviews.NewData true "UpdateReview Request"
// @Security BearerAuth
// @Success 200 {object} reviews.UpdateResp
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
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

	res, err := h.reviewClient.UpdateReview(ctx, &req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update review"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteReview godoc
// @Summary Delete review by ID
// @Description Delete a review by its ID
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Param("id")
	req := &pb.ID{Id: id}

	_, err := h.reviewClient.DeleteReview(ctx, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete review"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
