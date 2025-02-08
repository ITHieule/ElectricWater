package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterwaterRouter(router *gin.RouterGroup) {
	router.GET("/getwater", controllers.Water.Getwater)
	router.POST("/addwater", controllers.Water.AddWaterRecords)
	router.PUT("/updatewater", controllers.Water.UpdateWaterRecords)
	router.DELETE("/deletewater", controllers.Water.DeleteWaterRecords)
	router.POST("/searchwater", controllers.Water.SearchWaterRecords)
}
