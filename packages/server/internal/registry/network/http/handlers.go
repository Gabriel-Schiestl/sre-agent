package http

import (
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
	"github.com/gin-gonic/gin"
)

func listSuites(c *gin.Context, suiteService services.Svc[types.Suite]) {

}