package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterelectricRouter(router *gin.RouterGroup) {

	router.GET("/GetEnergyRecords", controllers.EnergyRecords.GetEnergyRecords)
	router.POST("/AddEnergyRecords", controllers.EnergyRecords.AddEnergyRecords)
	router.PUT("/UpdateEnergyRecords", controllers.EnergyRecords.UpdateEnergyRecords)
	router.DELETE("/DeleteEnergyRecords", controllers.EnergyRecords.DeleteEnergyRecords)
	router.POST("/SearchEnergyRecords", controllers.EnergyRecords.SearchEnergyRecords)
}
