package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type MicroserviceDB interface {
	List(ctx context.Context) []*types.Microservice
	ListBySuiteID(ctx context.Context, suiteID string) []*types.Microservice
	GetByID(ctx context.Context, id string) (*types.Microservice, error)
	Create(ctx context.Context, m *types.Microservice) (*types.Microservice, error)
	Update(ctx context.Context, m *types.Microservice) (*types.Microservice, error)
	Delete(ctx context.Context, id string) error
}

type microserviceDB struct {
	db *DB
}

func NewMicroserviceDB(db *DB) MicroserviceDB {
	return &microserviceDB{db: db}
}

const microserviceSelectCols = `
	id, test_suite_id, name, description, language, main_endpoints,
	cpu_limit, memory_limit, slo_latency_p99_ms, slo_error_rate_pct,
	prometheus_job_label, kubernetes_namespace, created_at`

func (r *microserviceDB) List(ctx context.Context) []*types.Microservice {
	rows, err := r.db.db.QueryContext(ctx,
		`SELECT`+microserviceSelectCols+` FROM microservices ORDER BY created_at DESC`,
	)
	if err != nil {
		log.Printf("microserviceDB.List: %v", err)
		return []*types.Microservice{}
	}
	defer rows.Close()

	var result []*types.Microservice
	for rows.Next() {
		m, err := scanMicroservice(rows)
		if err != nil {
			log.Printf("microserviceDB.List scan: %v", err)
			continue
		}
		result = append(result, m)
	}
	return result
}

func (r *microserviceDB) ListBySuiteID(ctx context.Context, suiteID string) []*types.Microservice {
	rows, err := r.db.db.QueryContext(ctx,
		`SELECT`+microserviceSelectCols+` FROM microservices WHERE test_suite_id = $1 ORDER BY created_at DESC`,
		suiteID,
	)
	if err != nil {
		log.Printf("microserviceDB.ListBySuiteID: %v", err)
		return []*types.Microservice{}
	}
	defer rows.Close()

	var result []*types.Microservice
	for rows.Next() {
		m, err := scanMicroservice(rows)
		if err != nil {
			log.Printf("microserviceDB.ListBySuiteID scan: %v", err)
			continue
		}
		result = append(result, m)
	}
	return result
}

func (r *microserviceDB) GetByID(ctx context.Context, id string) (*types.Microservice, error) {
	row := r.db.db.QueryRowContext(ctx,
		`SELECT`+microserviceSelectCols+` FROM microservices WHERE id = $1`,
		id,
	)
	return scanMicroserviceRow(row)
}

func (r *microserviceDB) Create(ctx context.Context, m *types.Microservice) (*types.Microservice, error) {
	endpointsJSON, err := json.Marshal(m.MainEndpoints())
	if err != nil {
		return nil, fmt.Errorf("microserviceDB.Create marshal endpoints: %w", err)
	}
	_, err = r.db.db.ExecContext(ctx, `
		INSERT INTO microservices
		(id, test_suite_id, name, description, language, main_endpoints,
		 cpu_limit, memory_limit, slo_latency_p99_ms, slo_error_rate_pct,
		 prometheus_job_label, kubernetes_namespace, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, m.ID(), m.TestSuiteID(), m.Name(), m.Description(), m.Language(),
		string(endpointsJSON), m.CPULimit(), m.MemoryLimit(),
		m.SLOLatencyP99Ms(), m.SLOErrorRatePct(),
		m.PrometheusJobLabel(), m.KubernetesNamespace(),
		m.CreatedAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("microserviceDB.Create: %w", err)
	}
	return m, nil
}

func (r *microserviceDB) Update(ctx context.Context, m *types.Microservice) (*types.Microservice, error) {
	endpointsJSON, err := json.Marshal(m.MainEndpoints())
	if err != nil {
		return nil, fmt.Errorf("microserviceDB.Update marshal endpoints: %w", err)
	}
	_, err = r.db.db.ExecContext(ctx, `
		UPDATE microservices
		SET name = $1, description = $2, language = $3, main_endpoints = $4,
		    cpu_limit = $5, memory_limit = $6, slo_latency_p99_ms = $7, slo_error_rate_pct = $8,
		    prometheus_job_label = $9, kubernetes_namespace = $10
		WHERE id = $11
	`, m.Name(), m.Description(), m.Language(), string(endpointsJSON),
		m.CPULimit(), m.MemoryLimit(), m.SLOLatencyP99Ms(), m.SLOErrorRatePct(),
		m.PrometheusJobLabel(), m.KubernetesNamespace(),
		m.ID(),
	)
	if err != nil {
		return nil, fmt.Errorf("microserviceDB.Update: %w", err)
	}
	return m, nil
}

func (r *microserviceDB) Delete(ctx context.Context, id string) error {
	_, err := r.db.db.ExecContext(ctx, `DELETE FROM microservices WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("microserviceDB.Delete: %w", err)
	}
	return nil
}

func scanMicroservice(rows *sql.Rows) (*types.Microservice, error) {
	var (
		id, suiteID, name, description, language, endpointsJSON, cpuLimit, memoryLimit string
		sloLatency                                                                       int
		sloErrorRate                                                                     float64
		prometheusJobLabel, kubernetesNamespace                                          sql.NullString
		createdAt                                                                        any
	)
	if err := rows.Scan(
		&id, &suiteID, &name, &description, &language, &endpointsJSON,
		&cpuLimit, &memoryLimit, &sloLatency, &sloErrorRate,
		&prometheusJobLabel, &kubernetesNamespace, &createdAt,
	); err != nil {
		return nil, err
	}
	return buildMicroservice(id, suiteID, name, description, language, endpointsJSON, cpuLimit, memoryLimit, sloLatency, sloErrorRate, prometheusJobLabel, kubernetesNamespace, createdAt)
}

func scanMicroserviceRow(row *sql.Row) (*types.Microservice, error) {
	var (
		id, suiteID, name, description, language, endpointsJSON, cpuLimit, memoryLimit string
		sloLatency                                                                       int
		sloErrorRate                                                                     float64
		prometheusJobLabel, kubernetesNamespace                                          sql.NullString
		createdAt                                                                        any
	)
	if err := row.Scan(
		&id, &suiteID, &name, &description, &language, &endpointsJSON,
		&cpuLimit, &memoryLimit, &sloLatency, &sloErrorRate,
		&prometheusJobLabel, &kubernetesNamespace, &createdAt,
	); err != nil {
		return nil, fmt.Errorf("microservice not found: %w", err)
	}
	return buildMicroservice(id, suiteID, name, description, language, endpointsJSON, cpuLimit, memoryLimit, sloLatency, sloErrorRate, prometheusJobLabel, kubernetesNamespace, createdAt)
}

func buildMicroservice(
	id, suiteID, name, description, language, endpointsJSON, cpuLimit, memoryLimit string,
	sloLatency int, sloErrorRate float64,
	prometheusJobLabel, kubernetesNamespace sql.NullString,
	createdAt any,
) (*types.Microservice, error) {
	var endpoints []string
	if err := json.Unmarshal([]byte(endpointsJSON), &endpoints); err != nil {
		endpoints = []string{}
	}

	var jobLabel *string
	if prometheusJobLabel.Valid {
		v := prometheusJobLabel.String
		jobLabel = &v
	}

	var k8sNamespace *string
	if kubernetesNamespace.Valid {
		v := kubernetesNamespace.String
		k8sNamespace = &v
	}

	return types.LoadMicroservice(
		id, suiteID, name, description, language, endpoints,
		cpuLimit, memoryLimit, sloLatency, sloErrorRate,
		jobLabel, k8sNamespace,
		toTime(createdAt),
	), nil
}
