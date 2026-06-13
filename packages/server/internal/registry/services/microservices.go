package services

import (
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type MicroserviceSvc interface {
	Svc[*types.Microservice]
	ListBySuiteID(suiteID string) []*types.Microservice
}

type microserviceService struct {
	db data.MicroserviceDB
}

func NewMicroserviceService(db data.MicroserviceDB) MicroserviceSvc {
	return &microserviceService{db: db}
}

func (s *microserviceService) List() []*types.Microservice {
	return s.db.List()
}

func (s *microserviceService) ListBySuiteID(suiteID string) []*types.Microservice {
	return s.db.ListBySuiteID(suiteID)
}

func (s *microserviceService) GetByID(id string) (*types.Microservice, error) {
	return s.db.GetByID(id)
}

func (s *microserviceService) Create(m *types.Microservice) (*types.Microservice, error) {
	return s.db.Create(m)
}

func (s *microserviceService) Update(id string, m *types.Microservice) (*types.Microservice, error) {
	existing, err := s.db.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("microservice not found: %w", err)
	}
	updated := types.LoadMicroservice(
		existing.ID(), existing.TestSuiteID(),
		m.Name(), m.Description(), m.Language(), m.MainEndpoints(),
		m.CPULimit(), m.MemoryLimit(), m.SLOLatencyP99Ms(), m.SLOErrorRatePct(),
		existing.CreatedAt(),
	)
	return s.db.Update(updated)
}

func (s *microserviceService) Delete(id string) error {
	if _, err := s.db.GetByID(id); err != nil {
		return fmt.Errorf("microservice not found: %w", err)
	}
	return s.db.Delete(id)
}
