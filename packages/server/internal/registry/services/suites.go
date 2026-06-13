package services

import (
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type Svc[T any] interface {
	List() []T
	GetByID(id string) (T, error)
	Delete(id string) error
	Create(item T) (T, error)
	Update(id string, item T) (T, error)
}

type suiteService struct{
	db data.SuiteDB
}

func NewSuiteService(db data.SuiteDB) Svc[types.Suite] {
	return &suiteService{db: db}
}

func (s *suiteService) List() []types.Suite {
	return []types.Suite{}
}

func (s *suiteService) GetByID(id string) (types.Suite, error) {
	return types.Suite{}, nil
}

func (s *suiteService) Delete(id string) error {
	return nil
}

func (s *suiteService) Create(item types.Suite) (types.Suite, error) {
	return types.Suite{}, nil
}

func (s *suiteService) Update(id string, item types.Suite) (types.Suite, error) {
	return types.Suite{}, nil
}