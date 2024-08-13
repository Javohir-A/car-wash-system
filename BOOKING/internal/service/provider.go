package service

import (
	"booking-service/genproto/providers"
	"booking-service/internal/storage"
	"context"
	"log/slog"
)

type ProviderService struct {
	log *slog.Logger
	providers.UnimplementedProvidersServer
	providerStorage storage.ProvidersStorage
}

func NewProviderService(log *slog.Logger, provider storage.ProvidersStorage) *ProviderService {
	return &ProviderService{
		log:             log,
		providerStorage: provider,
	}
}

func (p *ProviderService) CreateProvider(ctx context.Context, req *providers.NewProvider) (*providers.CreateResp, error) {
	res, err := p.providerStorage.CreateProvider(ctx, req)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	return res, nil
}

func (p *ProviderService) SearchProviders(ctx context.Context, req *providers.Filter) (*providers.SearchResp, error) {
	res, err := p.providerStorage.SearchProviders(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
