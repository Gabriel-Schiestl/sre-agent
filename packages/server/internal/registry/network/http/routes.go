package http

import (
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(suiteService services.Svc[types.Suite]) {
	suiteGroup := server.Group("/suite")

	suiteGroup.GET("", func(c *gin.Context) {
		listSuites(c, suiteService)
	})
}