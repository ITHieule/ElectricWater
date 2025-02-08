package controllers

import (
	"net/http"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type ElectricController struct {
	*BaseController
}

var EnergyRecords = &ElectricController{}

func (c *ElectricController) GetEnergyRecords(ctx *gin.Context) {
	result, err := services.Energy.IceLevelsSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
func (c *ElectricController) AddEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Energy.AddFlavorsSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

func (c *ElectricController) UpdateEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Energy.UpdateEnergySevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

func (c *ElectricController) DeleteEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	err := services.Energy.DeleteEnergySevice(requestParams.RecordID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, nil)
}
func (c *ElectricController) SearchEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Energy.SearchEnergySevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
