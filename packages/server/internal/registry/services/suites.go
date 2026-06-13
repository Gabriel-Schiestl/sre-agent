package services

import (
	"fmt"
	"time"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type Svc[T any] interface {
	GetByID(id string) (T, error)
	Delete(id string) error
	Create(item T) (T, error)
	Update(id string, item T) (T, error)
}

type SuiteSvc interface {
	Svc[*types.Suite]
	List() []*types.Suite
}

type suiteService struct {
	db data.SuiteDB
}

func NewSuiteService(db data.SuiteDB) SuiteSvc {
	return &suiteService{db: db}
}

func (s *suiteService) List() []*types.Suite {
	return s.db.List()
}

func (s *suiteService) GetByID(id string) (*types.Suite, error) {
	return s.db.GetByID(id)
}

func (s *suiteService) Delete(id string) error {
	if _, err := s.db.GetByID(id); err != nil {
		return fmt.Errorf("suite not found: %w", err)
	}
	return s.db.Delete(id)
}

func (s *suiteService) Create(item *types.Suite) (*types.Suite, error) {
	return s.db.Create(item)
}

func (s *suiteService) Update(id string, item *types.Suite) (*types.Suite, error) {
	existing, err := s.db.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("suite not found: %w", err)
	}
	updated := types.LoadSuite(existing.ID(), item.Name(), item.Description(), existing.CreatedAt(), time.Now())
	return s.db.Update(updated)
}
