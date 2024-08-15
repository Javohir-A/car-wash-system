package service

import (
	"booking-service/genproto/payments"
	"booking-service/internal/storage"
	"context"
	"errors"
	"log"
	"log/slog"
)

type PaymentsService struct {
	log *log.Logger
	payments.UnimplementedPaymentsServer
	paymentStorage storage.PaymentsStorage
}

func NewPaymentsService(log *log.Logger, paymentStorage storage.PaymentsStorage) *PaymentsService {
	return &PaymentsService{
		log:            log,
		paymentStorage: paymentStorage,
	}
}

func (s *PaymentsService) CreatePayment(ctx context.Context, req *payments.NewPayment) (*payments.CreateResp, error) {
	if req.BookingId == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		s.log.Println("validation failed for CreatePayment reason missing or invalid required fields")
		return nil, errors.New("missing or invalid required fields")
	}

	res, err := s.paymentStorage.CreatePayment(ctx, req)
	if err != nil {
		s.log.Println("failed to create payment", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *PaymentsService) GetPayment(ctx context.Context, req *payments.ID) (*payments.Payment, error) {
	if req.Id == "" {
		log.Println("validation failed for GetPayment")
		s.log.Println("validation failed for GetPayment", slog.String("reason", "missing payment ID"))
		return nil, errors.New("missing payment ID")
	}
	log.Println(req)

	payment, err := s.paymentStorage.GetPayment(ctx, req)
	if err != nil {
		log.Println(err)
		s.log.Println("failed to get payment", slog.Any("error", err))
		return nil, err
	}

	return payment, nil
}

func (s *PaymentsService) ListPayments(ctx context.Context, req *payments.Pagination) (*payments.PaymentsList, error) {
	res, err := s.paymentStorage.ListPayments(ctx, req)
	if err != nil {
		s.log.Println("failed to list payments", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
