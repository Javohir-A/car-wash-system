package service

import (
	"booking-service/genproto/bookings"
	"booking-service/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type BookingService struct {
	log *slog.Logger
	bookings.UnimplementedBookingsServer
	bookingStorage storage.BookingsStorage
}

func NewBookingService(log *slog.Logger, bookingStorage storage.BookingsStorage) *BookingService {
	return &BookingService{
		log:            log,
		bookingStorage: bookingStorage,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *bookings.NewBooking) (*bookings.CreateResp, error) {

	if req.UserId == "" || req.ProviderId == "" || req.ServiceId == "" {
		s.log.Error("validation failed for CreateBooking", slog.String("reason", "missing required fields"))
		return nil, errors.New("missing required fields")
	}

	res, err := s.bookingStorage.CreateBooking(ctx, req)
	if err != nil {
		s.log.Error("failed to create booking", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *bookings.NewData) (*bookings.UpdateResp, error) {

	if req.Id == "" {
		s.log.Error("validation failed for UpdateBooking", slog.String("reason", "missing booking ID"))
		return nil, errors.New("missing booking ID")
	}

	res, err := s.bookingStorage.UpdateBooking(ctx, req)
	if err != nil {
		s.log.Error("failed to update booking", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, req *bookings.ID) (*bookings.Void, error) {

	if req.Id == "" {
		s.log.Error("validation failed for CancelBooking", slog.String("reason", "missing booking ID"))
		return nil, errors.New("missing booking ID")
	}

	res, err := s.bookingStorage.CancelBooking(ctx, req)
	if err != nil {
		s.log.Error("failed to cancel booking", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *BookingService) ListBookings(ctx context.Context, req *bookings.Pagination) (*bookings.BookingsList, error) {
	res, err := s.bookingStorage.ListBookings(ctx, req)
	if err != nil {
		s.log.Error("failed to list bookings", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *bookings.ID) (*bookings.Booking, error) {

	if req.Id == "" {
		s.log.Error("validation failed for GetBooking", slog.String("reason", "missing booking ID"))
		return nil, errors.New("missing booking ID")
	}

	booking, err := s.bookingStorage.GetBooking(ctx, req)
	if err != nil {
		s.log.Error("failed to get booking", slog.Any("error", err))
		return nil, err
	}
	return booking, nil
}
