package service

import (
	"booking-service/genproto/payments"
	"booking-service/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type PaymentsService struct {
	log *slog.Logger
	payments.UnimplementedPaymentsServer
	paymentStorage storage.PaymentsStorage
}

func NewPaymentsService(log *slog.Logger, paymentStorage storage.PaymentsStorage) *PaymentsService {
	return &PaymentsService{
		log:            log,
		paymentStorage: paymentStorage,
	}
}

func (s *PaymentsService) CreatePayment(ctx context.Context, req *payments.NewPayment) (*payments.CreateResp, error) {
	if req.BookingId == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		s.log.Error("validation failed for CreatePayment", slog.String("reason", "missing or invalid required fields"))
		return nil, errors.New("missing or invalid required fields")
	}

	res, err := s.paymentStorage.CreatePayment(ctx, req)
	if err != nil {
		s.log.Error("failed to create payment", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *PaymentsService) GetPayment(ctx context.Context, req *payments.ID) (*payments.Payment, error) {
	if req.Id == "" {
		s.log.Error("validation failed for GetPayment", slog.String("reason", "missing payment ID"))
		return nil, errors.New("missing payment ID")
	}

	payment, err := s.paymentStorage.GetPayment(ctx, req)
	if err != nil {
		s.log.Error("failed to get payment", slog.Any("error", err))
		return nil, err
	}
	return payment, nil
}

func (s *PaymentsService) ListPayments(ctx context.Context, req *payments.Pagination) (*payments.PaymentsList, error) {
	res, err := s.paymentStorage.ListPayments(ctx, req)
	if err != nil {
		s.log.Error("failed to list payments", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
