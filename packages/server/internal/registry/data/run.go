package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type RunDB interface {
	List() []*types.TestRun
	ListBySuiteID(suiteID string) []*types.TestRun
	GetByID(id string) (*types.TestRun, error)
	Create(run *types.TestRun) (*types.TestRun, error)
	UpdateStatus(id string, status types.RunStatus) error
	Delete(id string) error
}

type runDB struct {
	db *DB
}

func NewRunDB(db *DB) RunDB {
	return &runDB{db: db}
}

func (r *runDB) List() []*types.TestRun {
	rows, err := r.db.db.Query(`
		SELECT id, test_suite_id, name, virtual_users, duration_seconds,
		       notes, status, jtl_file_path, created_at
		FROM test_runs ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("runDB.List: %v", err)
		return []*types.TestRun{}
	}
	defer rows.Close()

	var result []*types.TestRun
	for rows.Next() {
		run, err := scanRun(rows)
		if err != nil {
			log.Printf("runDB.List scan: %v", err)
			continue
		}
		result = append(result, run)
	}
	return result
}

func (r *runDB) ListBySuiteID(suiteID string) []*types.TestRun {
	rows, err := r.db.db.Query(`
		SELECT id, test_suite_id, name, virtual_users, duration_seconds,
		       notes, status, jtl_file_path, created_at
		FROM test_runs WHERE test_suite_id = $1 ORDER BY created_at DESC
	`, suiteID)
	if err != nil {
		log.Printf("runDB.ListBySuiteID: %v", err)
		return []*types.TestRun{}
	}
	defer rows.Close()

	var result []*types.TestRun
	for rows.Next() {
		run, err := scanRun(rows)
		if err != nil {
			log.Printf("runDB.ListBySuiteID scan: %v", err)
			continue
		}
		result = append(result, run)
	}
	return result
}

func (r *runDB) GetByID(id string) (*types.TestRun, error) {
	row := r.db.db.QueryRow(`
		SELECT id, test_suite_id, name, virtual_users, duration_seconds,
		       notes, status, jtl_file_path, created_at
		FROM test_runs WHERE id = $1
	`, id)
	return scanRunRow(row)
}

func (r *runDB) Create(run *types.TestRun) (*types.TestRun, error) {
	_, err := r.db.db.Exec(`
		INSERT INTO test_runs
		(id, test_suite_id, name, virtual_users, duration_seconds, notes, status, jtl_file_path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, run.ID(), run.TestSuiteID(), run.Name(), run.VirtualUsers(), run.DurationSeconds(),
		run.Notes(), string(run.Status()), run.JTLFilePath(), run.CreatedAt())
	if err != nil {
		return nil, fmt.Errorf("runDB.Create: %w", err)
	}
	return run, nil
}

func (r *runDB) UpdateStatus(id string, status types.RunStatus) error {
	_, err := r.db.db.Exec(`UPDATE test_runs SET status = $1 WHERE id = $2`, string(status), id)
	if err != nil {
		return fmt.Errorf("runDB.UpdateStatus: %w", err)
	}
	return nil
}

func (r *runDB) Delete(id string) error {
	_, err := r.db.db.Exec(`DELETE FROM test_runs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("runDB.Delete: %w", err)
	}
	return nil
}

func scanRun(rows *sql.Rows) (*types.TestRun, error) {
	var id, suiteID, name, notes, statusStr, jtlPath string
	var virtualUsers, durationSeconds int
	var createdAt any

	if err := rows.Scan(&id, &suiteID, &name, &virtualUsers, &durationSeconds, &notes, &statusStr, &jtlPath, &createdAt); err != nil {
		return nil, err
	}
	return types.LoadTestRun(id, suiteID, name, virtualUsers, durationSeconds, notes, types.RunStatus(statusStr), jtlPath, toTime(createdAt)), nil
}

func scanRunRow(row *sql.Row) (*types.TestRun, error) {
	var id, suiteID, name, notes, statusStr, jtlPath string
	var virtualUsers, durationSeconds int
	var createdAt any

	if err := row.Scan(&id, &suiteID, &name, &virtualUsers, &durationSeconds, &notes, &statusStr, &jtlPath, &createdAt); err != nil {
		return nil, fmt.Errorf("run not found: %w", err)
	}
	return types.LoadTestRun(id, suiteID, name, virtualUsers, durationSeconds, notes, types.RunStatus(statusStr), jtlPath, toTime(createdAt)), nil
}
