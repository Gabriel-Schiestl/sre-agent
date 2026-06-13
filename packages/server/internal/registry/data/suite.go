package data

import "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"

type SuiteDB interface {
	List() []types.Suite
	GetByID(id string) (types.Suite, error)
}

type suiteDB struct {
	db *DB
}

func NewSuiteDB(db *DB) SuiteDB {
	return &suiteDB{db: db}
}

func (s *suiteDB) List() []types.Suite {
	return []types.Suite{}
}

func (s *suiteDB) GetByID(id string) (types.Suite, error) {
	return types.Suite{}, nil
}