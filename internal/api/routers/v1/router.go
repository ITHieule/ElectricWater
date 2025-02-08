package router_v1

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	RegisterCommonRouter(v1.Group(""))

	RegisterReportRouter(v1.Group("/rp"))

	RegisterAdminRouter(v1.Group("/rp/admin"))

	RegisterAuthRouter(v1.Group("/auth"))

	RegisterFactoriesRouter(v1.Group("/factories"))

	RegisterElectricRouter(v1.Group("/energy"))

	RegisterWaterRouter(v1.Group("/water"))

}
