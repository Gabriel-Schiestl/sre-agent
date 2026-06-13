package types

import "time"

// --- Suite ---

type Suite struct {
	id          string
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewSuite(name, description string) *Suite {
	now := time.Now()
	return &Suite{
		id:          newID(),
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}
}

func LoadSuite(id, name, description string, createdAt, updatedAt time.Time) *Suite {
	return &Suite{id: id, name: name, description: description, createdAt: createdAt, updatedAt: updatedAt}
}

func (s *Suite) ID() string           { return s.id }
func (s *Suite) Name() string         { return s.name }
func (s *Suite) Description() string  { return s.description }
func (s *Suite) CreatedAt() time.Time { return s.createdAt }
func (s *Suite) UpdatedAt() time.Time { return s.updatedAt }

// --- Microservice ---

type Microservice struct {
	id              string
	testSuiteID     string
	name            string
	description     string
	language        string
	mainEndpoints   []string
	cpuLimit        string
	memoryLimit     string
	sloLatencyP99Ms int
	sloErrorRatePct float64
	createdAt       time.Time
}

func NewMicroservice(suiteID, name, description, language string, mainEndpoints []string, cpuLimit, memoryLimit string, sloLatencyP99Ms int, sloErrorRatePct float64) *Microservice {
	return &Microservice{
		id:              newID(),
		testSuiteID:     suiteID,
		name:            name,
		description:     description,
		language:        language,
		mainEndpoints:   mainEndpoints,
		cpuLimit:        cpuLimit,
		memoryLimit:     memoryLimit,
		sloLatencyP99Ms: sloLatencyP99Ms,
		sloErrorRatePct: sloErrorRatePct,
		createdAt:       time.Now(),
	}
}

func LoadMicroservice(id, testSuiteID, name, description, language string, mainEndpoints []string, cpuLimit, memoryLimit string, sloLatencyP99Ms int, sloErrorRatePct float64, createdAt time.Time) *Microservice {
	return &Microservice{
		id: id, testSuiteID: testSuiteID, name: name, description: description,
		language: language, mainEndpoints: mainEndpoints, cpuLimit: cpuLimit,
		memoryLimit: memoryLimit, sloLatencyP99Ms: sloLatencyP99Ms,
		sloErrorRatePct: sloErrorRatePct, createdAt: createdAt,
	}
}

func (m *Microservice) ID() string              { return m.id }
func (m *Microservice) TestSuiteID() string     { return m.testSuiteID }
func (m *Microservice) Name() string            { return m.name }
func (m *Microservice) Description() string     { return m.description }
func (m *Microservice) Language() string        { return m.language }
func (m *Microservice) MainEndpoints() []string { return m.mainEndpoints }
func (m *Microservice) CPULimit() string        { return m.cpuLimit }
func (m *Microservice) MemoryLimit() string     { return m.memoryLimit }
func (m *Microservice) SLOLatencyP99Ms() int    { return m.sloLatencyP99Ms }
func (m *Microservice) SLOErrorRatePct() float64 { return m.sloErrorRatePct }
func (m *Microservice) CreatedAt() time.Time    { return m.createdAt }

// --- TestRun ---

type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusAnalyzing RunStatus = "analyzing"
	RunStatusDone      RunStatus = "done"
	RunStatusFailed    RunStatus = "failed"
)

type TestRun struct {
	id              string
	testSuiteID     string
	name            string
	virtualUsers    int
	durationSeconds int
	notes           string
	status          RunStatus
	jtlFilePath     string
	createdAt       time.Time
}

func NewTestRun(suiteID, name string, virtualUsers, durationSeconds int, notes, jtlFilePath string) *TestRun {
	return &TestRun{
		id:              newID(),
		testSuiteID:     suiteID,
		name:            name,
		virtualUsers:    virtualUsers,
		durationSeconds: durationSeconds,
		notes:           notes,
		status:          RunStatusPending,
		jtlFilePath:     jtlFilePath,
		createdAt:       time.Now(),
	}
}

func LoadTestRun(id, testSuiteID, name string, virtualUsers, durationSeconds int, notes string, status RunStatus, jtlFilePath string, createdAt time.Time) *TestRun {
	return &TestRun{
		id: id, testSuiteID: testSuiteID, name: name, virtualUsers: virtualUsers,
		durationSeconds: durationSeconds, notes: notes, status: status,
		jtlFilePath: jtlFilePath, createdAt: createdAt,
	}
}

func (r *TestRun) ID() string           { return r.id }
func (r *TestRun) TestSuiteID() string  { return r.testSuiteID }
func (r *TestRun) Name() string         { return r.name }
func (r *TestRun) VirtualUsers() int    { return r.virtualUsers }
func (r *TestRun) DurationSeconds() int { return r.durationSeconds }
func (r *TestRun) Notes() string        { return r.notes }
func (r *TestRun) Status() RunStatus    { return r.status }
func (r *TestRun) JTLFilePath() string  { return r.jtlFilePath }
func (r *TestRun) CreatedAt() time.Time { return r.createdAt }

// --- Diagnosis ---

type ErrorCategory struct {
	Category          string   `json:"category"`
	Description       string   `json:"description"`
	Occurrences       int      `json:"occurrences"`
	AffectedEndpoints []string `json:"affectedEndpoints"`
	Severity          string   `json:"severity"` // low | medium | high | critical
}

type Hypothesis struct {
	Title    string `json:"title"`
	Evidence string `json:"evidence"`
	Priority int    `json:"priority"`
}

type Bottleneck struct {
	Microservice string       `json:"microservice"`
	Confidence   string       `json:"confidence"` // low | medium | high
	Hypotheses   []Hypothesis `json:"hypotheses"`
}

type Diagnosis struct {
	id          string
	testRunID   string
	errorPlan   []ErrorCategory
	bottlenecks []Bottleneck
	nextSteps   []string
	rawResponse string
	createdAt   time.Time
}

func NewDiagnosis(testRunID string, errorPlan []ErrorCategory, bottlenecks []Bottleneck, nextSteps []string, rawResponse string) *Diagnosis {
	return &Diagnosis{
		id:          newID(),
		testRunID:   testRunID,
		errorPlan:   errorPlan,
		bottlenecks: bottlenecks,
		nextSteps:   nextSteps,
		rawResponse: rawResponse,
		createdAt:   time.Now(),
	}
}

func LoadDiagnosis(id, testRunID string, errorPlan []ErrorCategory, bottlenecks []Bottleneck, nextSteps []string, rawResponse string, createdAt time.Time) *Diagnosis {
	return &Diagnosis{
		id: id, testRunID: testRunID, errorPlan: errorPlan, bottlenecks: bottlenecks,
		nextSteps: nextSteps, rawResponse: rawResponse, createdAt: createdAt,
	}
}

func (d *Diagnosis) ID() string                 { return d.id }
func (d *Diagnosis) TestRunID() string          { return d.testRunID }
func (d *Diagnosis) ErrorPlan() []ErrorCategory { return d.errorPlan }
func (d *Diagnosis) Bottlenecks() []Bottleneck  { return d.bottlenecks }
func (d *Diagnosis) NextSteps() []string        { return d.nextSteps }
func (d *Diagnosis) RawResponse() string        { return d.rawResponse }
func (d *Diagnosis) CreatedAt() time.Time       { return d.createdAt }
