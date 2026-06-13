package http

import "github.com/gin-gonic/gin"

var server *gin.Engine

func init() {
	server = gin.Default()
}

func Start(port string) error {
	return server.Run(":" + port)
}
