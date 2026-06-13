package data

import (
	"encoding/json"
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type DiagnosisDB interface {
	Save(d *types.Diagnosis) (*types.Diagnosis, error)
	GetByRunID(runID string) (*types.Diagnosis, error)
}

type diagnosisDB struct {
	db *DB
}

func NewDiagnosisDB(db *DB) DiagnosisDB {
	return &diagnosisDB{db: db}
}

func (r *diagnosisDB) Save(d *types.Diagnosis) (*types.Diagnosis, error) {
	errorPlanJSON, err := json.Marshal(d.ErrorPlan())
	if err != nil {
		return nil, fmt.Errorf("diagnosisDB.Save marshal error_plan: %w", err)
	}
	bottlenecksJSON, err := json.Marshal(d.Bottlenecks())
	if err != nil {
		return nil, fmt.Errorf("diagnosisDB.Save marshal bottlenecks: %w", err)
	}
	nextStepsJSON, err := json.Marshal(d.NextSteps())
	if err != nil {
		return nil, fmt.Errorf("diagnosisDB.Save marshal next_steps: %w", err)
	}

	_, err = r.db.db.Exec(`
		INSERT INTO diagnoses (id, test_run_id, error_plan, bottlenecks, next_steps, raw_response, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, d.ID(), d.TestRunID(), string(errorPlanJSON), string(bottlenecksJSON), string(nextStepsJSON), d.RawResponse(), d.CreatedAt())
	if err != nil {
		return nil, fmt.Errorf("diagnosisDB.Save: %w", err)
	}
	return d, nil
}

func (r *diagnosisDB) GetByRunID(runID string) (*types.Diagnosis, error) {
	row := r.db.db.QueryRow(`
		SELECT id, test_run_id, error_plan, bottlenecks, next_steps, raw_response, created_at
		FROM diagnoses WHERE test_run_id = $1
	`, runID)

	var id, testRunID, errorPlanJSON, bottlenecksJSON, nextStepsJSON, rawResponse string
	var createdAt any

	if err := row.Scan(&id, &testRunID, &errorPlanJSON, &bottlenecksJSON, &nextStepsJSON, &rawResponse, &createdAt); err != nil {
		return nil, fmt.Errorf("diagnosis not found for run %s: %w", runID, err)
	}

	var errorPlan []types.ErrorCategory
	if err := json.Unmarshal([]byte(errorPlanJSON), &errorPlan); err != nil {
		errorPlan = []types.ErrorCategory{}
	}

	var bottlenecks []types.Bottleneck
	if err := json.Unmarshal([]byte(bottlenecksJSON), &bottlenecks); err != nil {
		bottlenecks = []types.Bottleneck{}
	}

	var nextSteps []string
	if err := json.Unmarshal([]byte(nextStepsJSON), &nextSteps); err != nil {
		nextSteps = []string{}
	}

	return types.LoadDiagnosis(id, testRunID, errorPlan, bottlenecks, nextSteps, rawResponse, toTime(createdAt)), nil
}
