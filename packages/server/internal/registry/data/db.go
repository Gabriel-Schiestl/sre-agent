package data

import (
	"database/sql"
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/config"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func Open(cfg *config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS test_suites (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL,
			description TEXT NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS microservices (
			id                    TEXT PRIMARY KEY,
			test_suite_id         TEXT NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
			name                  TEXT NOT NULL,
			description           TEXT NOT NULL,
			language              TEXT NOT NULL,
			main_endpoints        TEXT NOT NULL DEFAULT '[]',
			cpu_limit             TEXT NOT NULL DEFAULT '',
			memory_limit          TEXT NOT NULL DEFAULT '',
			slo_latency_p99_ms    INTEGER NOT NULL DEFAULT 0,
			slo_error_rate_pct    REAL NOT NULL DEFAULT 0,
			prometheus_job_label  TEXT,
			kubernetes_namespace  TEXT,
			created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`ALTER TABLE microservices ADD COLUMN IF NOT EXISTS prometheus_job_label TEXT`,
		`ALTER TABLE microservices ADD COLUMN IF NOT EXISTS kubernetes_namespace TEXT`,
		`CREATE TABLE IF NOT EXISTS test_runs (
			id               TEXT PRIMARY KEY,
			test_suite_id    TEXT NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
			name             TEXT NOT NULL,
			virtual_users    INTEGER NOT NULL,
			duration_seconds INTEGER NOT NULL,
			notes            TEXT NOT NULL DEFAULT '',
			status           TEXT NOT NULL DEFAULT 'pending',
			jtl_file_path    TEXT NOT NULL DEFAULT '',
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS diagnoses (
			id           TEXT PRIMARY KEY,
			test_run_id  TEXT NOT NULL UNIQUE REFERENCES test_runs(id) ON DELETE CASCADE,
			error_plan   TEXT NOT NULL DEFAULT '[]',
			bottlenecks  TEXT NOT NULL DEFAULT '[]',
			next_steps   TEXT NOT NULL DEFAULT '[]',
			raw_response TEXT NOT NULL DEFAULT '',
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, q := range migrations {
		if _, err := d.db.Exec(q); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
