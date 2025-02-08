package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type FactoriesController struct {
	*BaseController
}

var Factories = &FactoriesController{}

func (c *FactoriesController) Getfactories(ctx *gin.Context) {
	result, err := services.Factor.FactoriesSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
