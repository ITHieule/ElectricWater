package router_v1

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	RegisterCommonRouter(v1.Group(""))

	RegisterReportRouter(v1.Group("/rp"))

	RegisterAdminRouter(v1.Group("/rp/admin"))

	RegisterelectricRouter(v1.Group("/energyrecords"))

	RegisterFactoriesRouter(v1.Group("/factories"))

	RegisterwaterRouter(v1.Group("/waterrecords"))

	RegisterAuthRouter(v1.Group("/auth"))

}
