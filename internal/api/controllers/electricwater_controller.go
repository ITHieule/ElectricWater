package controllers

import (
	"fmt"
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
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in IceLevelsSevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to get energy records: "+err.Error())
		return
	}

	// Kiểm tra nếu kết quả trả về rỗng
	if len(result) == 0 {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "No energy records found.")
		return
	}

	response.OkWithData(ctx, result)
}

func (c *ElectricController) AddEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Invalid request parameters:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid parameters: "+err.Error())
		return
	}

	result, err := services.Energy.AddEnergySevice(&requestParams)
	if err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in AddFlavorsSevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to add energy record: "+err.Error())
		return
	}

	// Kiểm tra kết quả trả về
	if result == nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to add energy record. No data returned.")
		return
	}

	response.OkWithData(ctx, result)
}

func (c *ElectricController) UpdateEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Invalid request parameters:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid parameters: "+err.Error())
		return
	}

	result, err := services.Energy.UpdateEnergySevice(&requestParams)
	if err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in UpdateEnergySevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to update energy record: "+err.Error())
		return
	}

	// Kiểm tra kết quả trả về
	if result == nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to update energy record. No data returned.")
		return
	}

	response.OkWithData(ctx, result)
}

func (c *ElectricController) DeleteEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Invalid request parameters:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid parameters: "+err.Error())
		return
	}

	err := services.Energy.DeleteEnergySevice(requestParams.RecordID)
	if err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in DeleteEnergySevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to delete energy record: "+err.Error())
		return
	}

	// Kiểm tra nếu RecordID không tồn tại
	if err == nil {
		response.OkWithData(ctx, nil)
	} else {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Energy record not found.")
	}
}

func (c *ElectricController) SearchEnergyRecords(ctx *gin.Context) {
	var requestParams request.Energyrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Invalid request parameters:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid parameters: "+err.Error())
		return
	}

	result, err := services.Energy.SearchEnergySevice(&requestParams)
	if err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in SearchEnergySevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to search energy records: "+err.Error())
		return
	}

	// Kiểm tra nếu không có kết quả tìm kiếm
	if len(result) == 0 {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "No energy records found matching the search criteria.")
		return
	}

	response.OkWithData(ctx, result)
}
