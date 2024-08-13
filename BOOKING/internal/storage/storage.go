package storage

import (
	"booking-service/genproto/providers"
	"booking-service/genproto/services"
	"context"
)

type ProvidersStorage interface {
	CreateProvider(context.Context, *providers.NewProvider) (*providers.CreateResp, error)
	SearchProviders(context.Context, *providers.Filter) (*providers.SearchResp, error)
}

type ServicesStorage interface {
	CreateService(context.Context, *services.NewService) (*services.CreateResp, error)
	DeleteService(context.Context, *services.ID) (*services.Void, error)
	ListServices(context.Context, *services.Pagination) (*services.ServicesList, error)
	SearchServices(context.Context, *services.Filter) (*services.SearchResp, error)
	UpdateService(context.Context, *services.NewData) (*services.UpdateResp, error)
	GetServiceByID(ctx context.Context, id *services.ID) (*services.Service, error) // New method
}
