package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	collectorpkg "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/collector"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type RunSvc interface {
	Svc[*types.TestRun]
	ListBySuiteID(ctx context.Context, suiteID string) []*types.TestRun
	CreateRun(ctx context.Context, run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) (*types.TestRun, error)
	GetDiagnosis(ctx context.Context, runID string) (*types.Diagnosis, error)
}

type runService struct {
	runDB       data.RunDB
	diagnosisDB data.DiagnosisDB
	runner      Runner
	analyst     Analyst
	collector   Collector
	uploadsDir  string
}

func NewRunService(runDB data.RunDB, diagnosisDB data.DiagnosisDB, runner Runner, analyst Analyst, collector Collector, uploadsDir string) RunSvc {
	return &runService{
		runDB:       runDB,
		diagnosisDB: diagnosisDB,
		runner:      runner,
		analyst:     analyst,
		collector:   collector,
		uploadsDir:  uploadsDir,
	}
}

// --- Svc[*types.TestRun] ---

func (s *runService) GetByID(ctx context.Context, id string) (*types.TestRun, error) {
	return s.runDB.GetByID(ctx, id)
}

// Create persists a run directly without triggering the analysis pipeline.
func (s *runService) Create(ctx context.Context, item *types.TestRun) (*types.TestRun, error) {
	return s.runDB.Create(ctx, item)
}

// Update is not supported — test runs are immutable after creation.
func (s *runService) Update(_ context.Context, _ string, _ *types.TestRun) (*types.TestRun, error) {
	return nil, fmt.Errorf("test runs are immutable")
}

func (s *runService) Delete(ctx context.Context, id string) error {
	run, err := s.runDB.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("run not found: %w", err)
	}
	if run.JTLFilePath() != "" {
		_ = os.Remove(run.JTLFilePath())
	}
	return s.runDB.Delete(ctx, id)
}

// --- RunSvc extensions ---

func (s *runService) ListBySuiteID(ctx context.Context, suiteID string) []*types.TestRun {
	return s.runDB.ListBySuiteID(ctx, suiteID)
}

func (s *runService) CreateRun(ctx context.Context, run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) (*types.TestRun, error) {
	jtlPath, err := s.saveJTL(run.ID(), jtlContent)
	if err != nil {
		return nil, fmt.Errorf("runService.CreateRun save jtl: %w", err)
	}

	persisted := types.NewTestRun(
		run.TestSuiteID(), run.Name(), run.VirtualUsers(),
		run.DurationSeconds(), run.Notes(), jtlPath,
	)

	created, err := s.runDB.Create(ctx, persisted)
	if err != nil {
		return nil, fmt.Errorf("runService.CreateRun persist: %w", err)
	}

	// Use background context: request context is cancelled after the response is sent.
	go s.process(context.Background(), created, suite, microservices, jtlContent)
	return created, nil
}

func (s *runService) GetDiagnosis(ctx context.Context, runID string) (*types.Diagnosis, error) {
	return s.diagnosisDB.GetByRunID(ctx, runID)
}

// --- internal ---

func (s *runService) saveJTL(runID string, content []byte) (string, error) {
	if err := os.MkdirAll(s.uploadsDir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(s.uploadsDir, runID+".jtl")
	if err := os.WriteFile(path, content, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (s *runService) process(ctx context.Context, run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) {
	if err := s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusAnalyzing); err != nil {
		log.Printf("runService.process update to analyzing: %v", err)
		return
	}

	if s.runner == nil || s.analyst == nil {
		log.Printf("runService.process: runner or analyst not configured for run %s", run.ID())
		_ = s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusFailed)
		return
	}

	aggregated, err := s.runner.Process(RunPayload{
		Run:           run,
		Suite:         suite,
		Microservices: microservices,
		JTLContent:    jtlContent,
	})
	if err != nil {
		log.Printf("runService.process runner failed for run %s: %v", run.ID(), err)
		_ = s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusFailed)
		return
	}

	var prometheusData *collectorpkg.PrometheusData
	if s.collector != nil {
		pd, err := s.collector.Collect(ctx, collectorpkg.CollectPayload{
			Microservices: microservices,
			StartTime:     aggregated.StartTime,
			EndTime:       aggregated.EndTime,
		})
		if err != nil {
			log.Printf("runService.process collector failed for run %s (continuing without metrics): %v", run.ID(), err)
		} else {
			prometheusData = pd
		}
	}

	diagnosis, err := s.analyst.Analyze(AnalysisPayload{
		Run:           run,
		Suite:         suite,
		Microservices: microservices,
		Data:          aggregated,
		Prometheus:    prometheusData,
	})
	if err != nil {
		log.Printf("runService.process analyst failed for run %s: %v", run.ID(), err)
		_ = s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusFailed)
		return
	}

	if _, err := s.diagnosisDB.Save(ctx, diagnosis); err != nil {
		log.Printf("runService.process save diagnosis for run %s: %v", run.ID(), err)
		_ = s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusFailed)
		return
	}

	if err := s.runDB.UpdateStatus(ctx, run.ID(), types.RunStatusDone); err != nil {
		log.Printf("runService.process update to done for run %s: %v", run.ID(), err)
	}
}
