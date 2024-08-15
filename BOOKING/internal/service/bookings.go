package service

import (
	"booking-service/genproto/bookings"
	"booking-service/internal/storage"
	"context"
	"errors"
	"log"
)

type BookingService struct {
	log *log.Logger
	bookings.UnimplementedBookingsServer
	bookingStorage storage.BookingsStorage
}

func NewBookingService(log *log.Logger, bookingStorage storage.BookingsStorage) *BookingService {
	return &BookingService{
		log:            log,
		bookingStorage: bookingStorage,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *bookings.NewBooking) (*bookings.CreateResp, error) {

	log.Println(req)
	if req.UserId == "" || req.ProviderId == "" || req.ServiceId == "" {
		s.log.Println("validation failed for CreateBooking: reason", "missing required fields")
		return nil, errors.New("missing required fields")
	}
	res, err := s.bookingStorage.CreateBooking(ctx, req)
	if err != nil {
		s.log.Println("failed to create booking error", err)
		return nil, err
	}

	return res, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *bookings.NewData) (*bookings.UpdateResp, error) {

	if req.Id == "" {
		s.log.Println("validation failed for UpdateBooking reason  missing booking ID")
		return nil, errors.New("missing booking ID")
	}

	res, err := s.bookingStorage.UpdateBooking(ctx, req)
	if err != nil {
		s.log.Println("failed to update booking error", err)
		return nil, err
	}
	return res, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, req *bookings.ID) (*bookings.Void, error) {

	if req.Id == "" {
		s.log.Println("validation failed for CancelBooking reason missing booking ID")
		return nil, errors.New("missing booking ID")
	}

	res, err := s.bookingStorage.CancelBooking(ctx, req)
	if err != nil {
		s.log.Println("failed to cancel bookingerror", err)
		return nil, err
	}
	return res, nil
}

func (s *BookingService) ListBookings(ctx context.Context, req *bookings.Pagination) (*bookings.BookingsList, error) {
	res, err := s.bookingStorage.ListBookings(ctx, req)
	if err != nil {
		s.log.Println("failed to list bookings error", err)
		return nil, err
	}
	return res, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *bookings.ID) (*bookings.Booking, error) {

	if req.Id == "" {
		s.log.Println("validation failed for GetBooking reason missing booking ID")
		return nil, errors.New("missing booking ID")
	}

	booking, err := s.bookingStorage.GetBooking(ctx, req)
	if err != nil {
		s.log.Println("failed to get booking error", err)
		return nil, err
	}
	return booking, nil
}
