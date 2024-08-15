package service

import (
	"booking-service/genproto/reviews"
	"booking-service/internal/storage"
	"context"
	"errors"
	"log"

	"log/slog"
)

type ReviewService struct {
	log *log.Logger
	reviews.UnimplementedReviewsServer
	reviewStorage storage.ReviewsStorage
}

func NewReviewService(log *log.Logger, reviewStorage storage.ReviewsStorage) *ReviewService {
	return &ReviewService{
		log:           log,
		reviewStorage: reviewStorage,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *reviews.NewReview) (*reviews.CreateResp, error) {
	if req.BookingId == "" || req.UserId == "" || req.ProviderId == "" || req.Rating == 0 {
		s.log.Println("validation failed for CreateReview", slog.String("reason", "missing required fields"))
		return nil, errors.New("missing required fields")
	}

	res, err := s.reviewStorage.CreateReview(ctx, req)
	if err != nil {
		s.log.Println("failed to create review", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *reviews.NewData) (*reviews.UpdateResp, error) {
	if req.Id == "" {
		s.log.Println("validation failed for UpdateReview", slog.String("reason", "missing review ID"))
		return nil, errors.New("missing review ID")
	}

	res, err := s.reviewStorage.UpdateReview(ctx, req)
	if err != nil {
		s.log.Println("failed to update review", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *reviews.ID) (*reviews.Void, error) {
	if req.Id == "" {
		s.log.Println("validation failed for DeleteReview", slog.String("reason", "missing review ID"))
		return nil, errors.New("missing review ID")
	}

	res, err := s.reviewStorage.DeleteReview(ctx, req)
	if err != nil {
		s.log.Println("failed to delete review", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ReviewService) ListReviews(ctx context.Context, req *reviews.Pagination) (*reviews.ReviewsList, error) {
	res, err := s.reviewStorage.ListReviews(ctx, req)
	if err != nil {
		s.log.Println("failed to list reviews", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
