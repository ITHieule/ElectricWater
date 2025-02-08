package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFactoriesRouter(router *gin.RouterGroup) {
	router.GET("/Getfactories", controllers.Factories.Getfactories)
}
