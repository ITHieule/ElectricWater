package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterWaterRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Water.Getwater)
	router.POST("/add", controllers.Water.AddWaterRecords)
	router.PUT("/update", controllers.Water.UpdateWaterRecords)
	router.DELETE("/delete", controllers.Water.DeleteWaterRecords)
	router.POST("/search", controllers.Water.SearchWaterRecords)
}
