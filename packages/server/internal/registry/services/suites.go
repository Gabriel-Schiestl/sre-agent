package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type Svc[T any] interface {
	GetByID(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, item T) (T, error)
	Update(ctx context.Context, id string, item T) (T, error)
}

type SuiteSvc interface {
	Svc[*types.Suite]
	List(ctx context.Context) []*types.Suite
}

type suiteService struct {
	db data.SuiteDB
}

func NewSuiteService(db data.SuiteDB) SuiteSvc {
	return &suiteService{db: db}
}

func (s *suiteService) List(ctx context.Context) []*types.Suite {
	return s.db.List(ctx)
}

func (s *suiteService) GetByID(ctx context.Context, id string) (*types.Suite, error) {
	return s.db.GetByID(ctx, id)
}

func (s *suiteService) Delete(ctx context.Context, id string) error {
	if _, err := s.db.GetByID(ctx, id); err != nil {
		return fmt.Errorf("suite not found: %w", err)
	}
	return s.db.Delete(ctx, id)
}

func (s *suiteService) Create(ctx context.Context, item *types.Suite) (*types.Suite, error) {
	return s.db.Create(ctx, item)
}

func (s *suiteService) Update(ctx context.Context, id string, item *types.Suite) (*types.Suite, error) {
	existing, err := s.db.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("suite not found: %w", err)
	}
	updated := types.LoadSuite(existing.ID(), item.Name(), item.Description(), existing.CreatedAt(), time.Now())
	return s.db.Update(ctx, updated)
}
