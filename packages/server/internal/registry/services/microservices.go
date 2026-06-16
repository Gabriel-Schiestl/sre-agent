package services

import (
	"context"
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type MicroserviceSvc interface {
	Svc[*types.Microservice]
	ListBySuiteID(ctx context.Context, suiteID string) []*types.Microservice
}

type microserviceService struct {
	db data.MicroserviceDB
}

func NewMicroserviceService(db data.MicroserviceDB) MicroserviceSvc {
	return &microserviceService{db: db}
}

func (s *microserviceService) List(ctx context.Context) []*types.Microservice {
	return s.db.List(ctx)
}

func (s *microserviceService) ListBySuiteID(ctx context.Context, suiteID string) []*types.Microservice {
	return s.db.ListBySuiteID(ctx, suiteID)
}

func (s *microserviceService) GetByID(ctx context.Context, id string) (*types.Microservice, error) {
	return s.db.GetByID(ctx, id)
}

func (s *microserviceService) Create(ctx context.Context, m *types.Microservice) (*types.Microservice, error) {
	return s.db.Create(ctx, m)
}

func (s *microserviceService) Update(ctx context.Context, id string, m *types.Microservice) (*types.Microservice, error) {
	existing, err := s.db.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("microservice not found: %w", err)
	}
	updated := types.LoadMicroservice(
		existing.ID(), existing.TestSuiteID(),
		m.Name(), m.Description(), m.Language(), m.MainEndpoints(),
		m.CPULimit(), m.MemoryLimit(), m.SLOLatencyP99Ms(), m.SLOErrorRatePct(),
		m.PrometheusJobLabel(), m.KubernetesNamespace(),
		existing.CreatedAt(),
	)
	return s.db.Update(ctx, updated)
}

func (s *microserviceService) Delete(ctx context.Context, id string) error {
	if _, err := s.db.GetByID(ctx, id); err != nil {
		return fmt.Errorf("microservice not found: %w", err)
	}
	return s.db.Delete(ctx, id)
}
