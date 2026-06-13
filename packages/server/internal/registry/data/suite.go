package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type SuiteDB interface {
	List(ctx context.Context) []*types.Suite
	GetByID(ctx context.Context, id string) (*types.Suite, error)
	Create(ctx context.Context, suite *types.Suite) (*types.Suite, error)
	Update(ctx context.Context, suite *types.Suite) (*types.Suite, error)
	Delete(ctx context.Context, id string) error
}

type suiteDB struct {
	db *DB
}

func NewSuiteDB(db *DB) SuiteDB {
	return &suiteDB{db: db}
}

func (s *suiteDB) List(ctx context.Context) []*types.Suite {
	rows, err := s.db.db.QueryContext(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM test_suites ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("suiteDB.List: %v", err)
		return []*types.Suite{}
	}
	defer rows.Close()

	var suites []*types.Suite
	for rows.Next() {
		suite, err := scanSuite(rows)
		if err != nil {
			log.Printf("suiteDB.List scan: %v", err)
			continue
		}
		suites = append(suites, suite)
	}
	return suites
}

func (s *suiteDB) GetByID(ctx context.Context, id string) (*types.Suite, error) {
	row := s.db.db.QueryRowContext(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM test_suites WHERE id = $1
	`, id)
	return scanSuiteRow(row)
}

func (s *suiteDB) Create(ctx context.Context, suite *types.Suite) (*types.Suite, error) {
	_, err := s.db.db.ExecContext(ctx, `
		INSERT INTO test_suites (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, suite.ID(), suite.Name(), suite.Description(), suite.CreatedAt(), suite.UpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("suiteDB.Create: %w", err)
	}
	return suite, nil
}

func (s *suiteDB) Update(ctx context.Context, suite *types.Suite) (*types.Suite, error) {
	_, err := s.db.db.ExecContext(ctx, `
		UPDATE test_suites SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`, suite.Name(), suite.Description(), suite.UpdatedAt(), suite.ID())
	if err != nil {
		return nil, fmt.Errorf("suiteDB.Update: %w", err)
	}
	return suite, nil
}

func (s *suiteDB) Delete(ctx context.Context, id string) error {
	_, err := s.db.db.ExecContext(ctx, `DELETE FROM test_suites WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("suiteDB.Delete: %w", err)
	}
	return nil
}

func scanSuite(rows *sql.Rows) (*types.Suite, error) {
	var id, name, description string
	var createdAt, updatedAt any
	if err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return types.LoadSuite(id, name, description, toTime(createdAt), toTime(updatedAt)), nil
}

func scanSuiteRow(row *sql.Row) (*types.Suite, error) {
	var id, name, description string
	var createdAt, updatedAt any
	if err := row.Scan(&id, &name, &description, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("suite not found: %w", err)
	}
	return types.LoadSuite(id, name, description, toTime(createdAt), toTime(updatedAt)), nil
}
