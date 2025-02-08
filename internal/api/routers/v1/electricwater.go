package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterelectricRouter(router *gin.RouterGroup) {

	router.GET("/get", controllers.EnergyRecords.GetEnergyRecords)
	router.POST("/add", controllers.EnergyRecords.AddEnergyRecords)
	router.PUT("/update", controllers.EnergyRecords.UpdateEnergyRecords)
	router.DELETE("/delete", controllers.EnergyRecords.DeleteEnergyRecords)
	router.POST("/search", controllers.EnergyRecords.SearchEnergyRecords)
}
