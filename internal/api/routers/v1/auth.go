package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRouter(router *gin.RouterGroup) {
	router.POST("/login", controllers.Login)
}
