package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type RunSvc interface {
	Svc[*types.TestRun]
	ListBySuiteID(suiteID string) []*types.TestRun
	CreateRun(run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) (*types.TestRun, error)
	GetDiagnosis(runID string) (*types.Diagnosis, error)
}

type runService struct {
	runDB       data.RunDB
	diagnosisDB data.DiagnosisDB
	runner      Runner
	analyst     Analyst
	uploadsDir  string
}

func NewRunService(runDB data.RunDB, diagnosisDB data.DiagnosisDB, runner Runner, analyst Analyst, uploadsDir string) RunSvc {
	return &runService{
		runDB:       runDB,
		diagnosisDB: diagnosisDB,
		runner:      runner,
		analyst:     analyst,
		uploadsDir:  uploadsDir,
	}
}

// --- Svc[*types.TestRun] ---

func (s *runService) GetByID(id string) (*types.TestRun, error) {
	return s.runDB.GetByID(id)
}

// Create persists a run directly without triggering the analysis pipeline.
func (s *runService) Create(item *types.TestRun) (*types.TestRun, error) {
	return s.runDB.Create(item)
}

// Update is not supported — test runs are immutable after creation.
func (s *runService) Update(_ string, _ *types.TestRun) (*types.TestRun, error) {
	return nil, fmt.Errorf("test runs are immutable")
}

func (s *runService) Delete(id string) error {
	run, err := s.runDB.GetByID(id)
	if err != nil {
		return fmt.Errorf("run not found: %w", err)
	}
	if run.JTLFilePath() != "" {
		_ = os.Remove(run.JTLFilePath())
	}
	return s.runDB.Delete(id)
}

// --- RunSvc extensions ---

func (s *runService) ListBySuiteID(suiteID string) []*types.TestRun {
	return s.runDB.ListBySuiteID(suiteID)
}

func (s *runService) CreateRun(run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) (*types.TestRun, error) {
	jtlPath, err := s.saveJTL(run.ID(), jtlContent)
	if err != nil {
		return nil, fmt.Errorf("runService.CreateRun save jtl: %w", err)
	}

	persisted := types.NewTestRun(
		run.TestSuiteID(), run.Name(), run.VirtualUsers(),
		run.DurationSeconds(), run.Notes(), jtlPath,
	)

	created, err := s.runDB.Create(persisted)
	if err != nil {
		return nil, fmt.Errorf("runService.CreateRun persist: %w", err)
	}

	go s.process(created, suite, microservices, jtlContent)
	return created, nil
}

func (s *runService) GetDiagnosis(runID string) (*types.Diagnosis, error) {
	return s.diagnosisDB.GetByRunID(runID)
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

func (s *runService) process(run *types.TestRun, suite *types.Suite, microservices []*types.Microservice, jtlContent []byte) {
	if err := s.runDB.UpdateStatus(run.ID(), types.RunStatusAnalyzing); err != nil {
		log.Printf("runService.process update to analyzing: %v", err)
		return
	}

	if s.runner == nil || s.analyst == nil {
		log.Printf("runService.process: runner or analyst not configured for run %s", run.ID())
		_ = s.runDB.UpdateStatus(run.ID(), types.RunStatusFailed)
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
		_ = s.runDB.UpdateStatus(run.ID(), types.RunStatusFailed)
		return
	}

	diagnosis, err := s.analyst.Analyze(AnalysisPayload{
		Run:           run,
		Suite:         suite,
		Microservices: microservices,
		Data:          aggregated,
	})
	if err != nil {
		log.Printf("runService.process analyst failed for run %s: %v", run.ID(), err)
		_ = s.runDB.UpdateStatus(run.ID(), types.RunStatusFailed)
		return
	}

	if _, err := s.diagnosisDB.Save(diagnosis); err != nil {
		log.Printf("runService.process save diagnosis for run %s: %v", run.ID(), err)
		_ = s.runDB.UpdateStatus(run.ID(), types.RunStatusFailed)
		return
	}

	if err := s.runDB.UpdateStatus(run.ID(), types.RunStatusDone); err != nil {
		log.Printf("runService.process update to done for run %s: %v", run.ID(), err)
	}
}
