package service

import (
	"booking-service/genproto/services"
	"booking-service/internal/storage"
	"context"
	"log/slog"
)

type ServicesService struct {
	log *slog.Logger
	services.UnimplementedServicesServer
	serviceStorage storage.ServicesStorage
}

func NewServicesService(log *slog.Logger, serviceStorage storage.ServicesStorage) *ServicesService {
	return &ServicesService{
		log:            log,
		serviceStorage: serviceStorage,
	}
}

func (s *ServicesService) CreateService(ctx context.Context, req *services.NewService) (*services.CreateResp, error) {
	res, err := s.serviceStorage.CreateService(ctx, req)
	if err != nil {
		s.log.Error("failed to create service", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ServicesService) UpdateService(ctx context.Context, req *services.NewData) (*services.UpdateResp, error) {
	res, err := s.serviceStorage.UpdateService(ctx, req)
	if err != nil {
		s.log.Error("failed to update service", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ServicesService) DeleteService(ctx context.Context, req *services.ID) (*services.Void, error) {
	res, err := s.serviceStorage.DeleteService(ctx, req)
	if err != nil {
		s.log.Error("failed to delete service", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ServicesService) ListServices(ctx context.Context, req *services.Pagination) (*services.ServicesList, error) {
	res, err := s.serviceStorage.ListServices(ctx, req)
	if err != nil {
		s.log.Error("failed to list services", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ServicesService) SearchServices(ctx context.Context, req *services.Filter) (*services.SearchResp, error) {
	res, err := s.serviceStorage.SearchServices(ctx, req)
	if err != nil {
		s.log.Error("failed to search services", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}

func (s *ServicesService) GetService(ctx context.Context, req *services.ID) (*services.Service, error) {
	service, err := s.serviceStorage.GetServiceByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return service, nil
}
