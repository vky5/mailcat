package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vky5/mailcat/internal/api/routes"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	// register all routes
	routes.RegisterAccountRoutes(r)

	return r
}
