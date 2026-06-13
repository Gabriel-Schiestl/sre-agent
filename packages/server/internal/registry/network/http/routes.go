package http

import (
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(suites services.SuiteSvc, microservices services.MicroserviceSvc, runs services.RunSvc) {
	h := NewHandlers(suites, microservices, runs)

	suiteGroup := server.Group("/suites")
	{
		suiteGroup.GET("", h.listSuites)
		suiteGroup.POST("", h.createSuite)
		suiteGroup.GET("/:id", h.getSuite)
		suiteGroup.PUT("/:id", h.updateSuite)
		suiteGroup.DELETE("/:id", h.deleteSuite)

		suiteGroup.POST("/:id/microservices", h.createMicroservice)
		suiteGroup.GET("/:id/runs", h.listRuns)
		suiteGroup.POST("/:id/runs", h.createRun)
	}

	microserviceGroup := server.Group("/microservices")
	{
		microserviceGroup.PUT("/:id", h.updateMicroservice)
		microserviceGroup.DELETE("/:id", h.deleteMicroservice)
	}

	runGroup := server.Group("/runs")
	{
		runGroup.GET("/:id", h.getRun)
		runGroup.GET("/:id/diagnosis", h.getDiagnosis)
	}
}

func SetupCORS(allowedOrigin string) {
	server.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}
