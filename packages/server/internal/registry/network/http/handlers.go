package http

import (
	"io"
	"net/http"
	"time"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
	"github.com/gin-gonic/gin"
)

// --- Handlers struct ---

type Handlers struct {
	suites        services.SuiteSvc
	microservices services.MicroserviceSvc
	runs          services.RunSvc
}

func NewHandlers(suites services.SuiteSvc, microservices services.MicroserviceSvc, runs services.RunSvc) *Handlers {
	return &Handlers{suites: suites, microservices: microservices, runs: runs}
}

// --- Request DTOs ---

type createSuiteRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required"`
}

type updateSuiteRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required"`
}

type createMicroserviceRequest struct {
	Name                string   `json:"name" binding:"required"`
	Description         string   `json:"description" binding:"required"`
	Language            string   `json:"language" binding:"required"`
	MainEndpoints       []string `json:"mainEndpoints"`
	CPULimit            string   `json:"cpuLimit"`
	MemoryLimit         string   `json:"memoryLimit"`
	SLOLatencyP99Ms     int      `json:"sloLatencyP99Ms"`
	SLOErrorRatePct     float64  `json:"sloErrorRatePct"`
	PrometheusJobLabel  *string  `json:"prometheusJobLabel"`
	KubernetesNamespace *string  `json:"kubernetesNamespace"`
}

type updateMicroserviceRequest struct {
	Name                string   `json:"name" binding:"required"`
	Description         string   `json:"description" binding:"required"`
	Language            string   `json:"language" binding:"required"`
	MainEndpoints       []string `json:"mainEndpoints"`
	CPULimit            string   `json:"cpuLimit"`
	MemoryLimit         string   `json:"memoryLimit"`
	SLOLatencyP99Ms     int      `json:"sloLatencyP99Ms"`
	SLOErrorRatePct     float64  `json:"sloErrorRatePct"`
	PrometheusJobLabel  *string  `json:"prometheusJobLabel"`
	KubernetesNamespace *string  `json:"kubernetesNamespace"`
}

type createRunRequest struct {
	Name            string `form:"name" binding:"required"`
	VirtualUsers    int    `form:"virtualUsers" binding:"required,min=1"`
	DurationSeconds int    `form:"durationSeconds" binding:"required,min=1"`
	Notes           string `form:"notes"`
}

// --- Response DTOs ---

type suiteResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type suiteDetailResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
	Microservices []microserviceResponse `json:"microservices"`
}

type microserviceResponse struct {
	ID                  string   `json:"id"`
	TestSuiteID         string   `json:"testSuiteId"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Language            string   `json:"language"`
	MainEndpoints       []string `json:"mainEndpoints"`
	CPULimit            string   `json:"cpuLimit"`
	MemoryLimit         string   `json:"memoryLimit"`
	SLOLatencyP99Ms     int      `json:"sloLatencyP99Ms"`
	SLOErrorRatePct     float64  `json:"sloErrorRatePct"`
	PrometheusJobLabel  *string  `json:"prometheusJobLabel"`
	KubernetesNamespace *string  `json:"kubernetesNamespace"`
	CreatedAt           string   `json:"createdAt"`
}

type runResponse struct {
	ID              string `json:"id"`
	TestSuiteID     string `json:"testSuiteId"`
	Name            string `json:"name"`
	VirtualUsers    int    `json:"virtualUsers"`
	DurationSeconds int    `json:"durationSeconds"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
}

type diagnosisResponse struct {
	ID          string                `json:"id"`
	TestRunID   string                `json:"testRunId"`
	ErrorPlan   []types.ErrorCategory `json:"errorPlan"`
	Bottlenecks []types.Bottleneck    `json:"bottlenecks"`
	NextSteps   []string              `json:"nextSteps"`
	CreatedAt   string                `json:"createdAt"`
}

// --- Converters ---

func toSuiteResponse(s *types.Suite) suiteResponse {
	return suiteResponse{
		ID:          s.ID(),
		Name:        s.Name(),
		Description: s.Description(),
		CreatedAt:   s.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt().Format(time.RFC3339),
	}
}

func toSuiteDetailResponse(s *types.Suite, ms []*types.Microservice) suiteDetailResponse {
	microservicesResp := make([]microserviceResponse, 0, len(ms))
	for _, m := range ms {
		microservicesResp = append(microservicesResp, toMicroserviceResponse(m))
	}
	return suiteDetailResponse{
		ID:            s.ID(),
		Name:          s.Name(),
		Description:   s.Description(),
		CreatedAt:     s.CreatedAt().Format(time.RFC3339),
		UpdatedAt:     s.UpdatedAt().Format(time.RFC3339),
		Microservices: microservicesResp,
	}
}

func toMicroserviceResponse(m *types.Microservice) microserviceResponse {
	endpoints := m.MainEndpoints()
	if endpoints == nil {
		endpoints = []string{}
	}
	return microserviceResponse{
		ID:                  m.ID(),
		TestSuiteID:         m.TestSuiteID(),
		Name:                m.Name(),
		Description:         m.Description(),
		Language:            m.Language(),
		MainEndpoints:       endpoints,
		CPULimit:            m.CPULimit(),
		MemoryLimit:         m.MemoryLimit(),
		SLOLatencyP99Ms:     m.SLOLatencyP99Ms(),
		SLOErrorRatePct:     m.SLOErrorRatePct(),
		PrometheusJobLabel:  m.PrometheusJobLabel(),
		KubernetesNamespace: m.KubernetesNamespace(),
		CreatedAt:           m.CreatedAt().Format(time.RFC3339),
	}
}

func toRunResponse(r *types.TestRun) runResponse {
	return runResponse{
		ID:              r.ID(),
		TestSuiteID:     r.TestSuiteID(),
		Name:            r.Name(),
		VirtualUsers:    r.VirtualUsers(),
		DurationSeconds: r.DurationSeconds(),
		Notes:           r.Notes(),
		Status:          string(r.Status()),
		CreatedAt:       r.CreatedAt().Format(time.RFC3339),
	}
}

func toDiagnosisResponse(d *types.Diagnosis) diagnosisResponse {
	errorPlan := d.ErrorPlan()
	if errorPlan == nil {
		errorPlan = []types.ErrorCategory{}
	}
	bottlenecks := d.Bottlenecks()
	if bottlenecks == nil {
		bottlenecks = []types.Bottleneck{}
	}
	nextSteps := d.NextSteps()
	if nextSteps == nil {
		nextSteps = []string{}
	}
	return diagnosisResponse{
		ID:          d.ID(),
		TestRunID:   d.TestRunID(),
		ErrorPlan:   errorPlan,
		Bottlenecks: bottlenecks,
		NextSteps:   nextSteps,
		CreatedAt:   d.CreatedAt().Format(time.RFC3339),
	}
}

// --- Suite handlers ---

func (h *Handlers) listSuites(c *gin.Context) {
	ctx := c.Request.Context()
	suites := h.suites.List(ctx)
	resp := make([]suiteResponse, 0, len(suites))
	for _, s := range suites {
		resp = append(resp, toSuiteResponse(s))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) getSuite(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	suite, err := h.suites.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "suite not found"})
		return
	}
	microservices := h.microservices.ListBySuiteID(ctx, id)
	c.JSON(http.StatusOK, toSuiteDetailResponse(suite, microservices))
}

func (h *Handlers) createSuite(c *gin.Context) {
	ctx := c.Request.Context()
	var req createSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	suite, err := h.suites.Create(ctx, types.NewSuite(req.Name, req.Description))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create suite"})
		return
	}
	c.JSON(http.StatusCreated, toSuiteResponse(suite))
}

func (h *Handlers) updateSuite(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var req updateSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.suites.Update(ctx, id, types.NewSuite(req.Name, req.Description))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toSuiteResponse(updated))
}

func (h *Handlers) deleteSuite(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.suites.Delete(ctx, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Microservice handlers ---

func (h *Handlers) createMicroservice(c *gin.Context) {
	ctx := c.Request.Context()
	suiteID := c.Param("id")
	var req createMicroserviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := h.suites.GetByID(ctx, suiteID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "suite not found"})
		return
	}
	endpoints := req.MainEndpoints
	if endpoints == nil {
		endpoints = []string{}
	}
	m := types.NewMicroservice(suiteID, req.Name, req.Description, req.Language, endpoints, req.CPULimit, req.MemoryLimit, req.SLOLatencyP99Ms, req.SLOErrorRatePct, req.PrometheusJobLabel, req.KubernetesNamespace)
	created, err := h.microservices.Create(ctx, m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create microservice"})
		return
	}
	c.JSON(http.StatusCreated, toMicroserviceResponse(created))
}

func (h *Handlers) updateMicroservice(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var req updateMicroserviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endpoints := req.MainEndpoints
	if endpoints == nil {
		endpoints = []string{}
	}
	m := types.NewMicroservice("", req.Name, req.Description, req.Language, endpoints, req.CPULimit, req.MemoryLimit, req.SLOLatencyP99Ms, req.SLOErrorRatePct, req.PrometheusJobLabel, req.KubernetesNamespace)
	updated, err := h.microservices.Update(ctx, id, m)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toMicroserviceResponse(updated))
}

func (h *Handlers) deleteMicroservice(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.microservices.Delete(ctx, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Run handlers ---

func (h *Handlers) listRuns(c *gin.Context) {
	ctx := c.Request.Context()
	suiteID := c.Param("id")
	runs := h.runs.ListBySuiteID(ctx, suiteID)
	resp := make([]runResponse, 0, len(runs))
	for _, r := range runs {
		resp = append(resp, toRunResponse(r))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) getRun(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	run, err := h.runs.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "run not found"})
		return
	}
	c.JSON(http.StatusOK, toRunResponse(run))
}

func (h *Handlers) createRun(c *gin.Context) {
	ctx := c.Request.Context()
	suiteID := c.Param("id")

	suite, err := h.suites.GetByID(ctx, suiteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "suite not found"})
		return
	}

	var req createRunRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("jtlFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jtlFile is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}
	defer file.Close()

	jtlContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}

	microservices := h.microservices.ListBySuiteID(ctx, suiteID)
	run := types.NewTestRun(suiteID, req.Name, req.VirtualUsers, req.DurationSeconds, req.Notes, "")

	created, err := h.runs.CreateRun(ctx, run, suite, microservices, jtlContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create run"})
		return
	}
	c.JSON(http.StatusCreated, toRunResponse(created))
}

func (h *Handlers) getDiagnosis(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("id")
	run, err := h.runs.GetByID(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "run not found"})
		return
	}
	if run.Status() != types.RunStatusDone {
		c.JSON(http.StatusAccepted, gin.H{"status": string(run.Status())})
		return
	}
	diagnosis, err := h.runs.GetDiagnosis(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "diagnosis not found"})
		return
	}
	c.JSON(http.StatusOK, toDiagnosisResponse(diagnosis))
}
