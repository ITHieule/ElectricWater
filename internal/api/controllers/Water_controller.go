package controllers

import (
	"errors"
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type WaterController struct {
	*BaseController
}

var Water = &WaterController{}

func (c *WaterController) Getwater(ctx *gin.Context) {
	result, err := services.WaterRecord.WaterRecordSevice()
	if err != nil {
		// Bắt lỗi và trả về thông báo chi tiết
		c.handleError(ctx, err)
		return
	}
	response.OkWithData(ctx, result)
}

func (c *WaterController) AddWaterRecords(ctx *gin.Context) {
	var requestParams request.Waterrequest

	// Kiểm tra và validate tham số yêu cầu
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		c.handleError(ctx, err)
		return
	}
	result, err := services.WaterRecord.AddWaterRecordSevice(&requestParams)
	if err != nil {
		// Bắt lỗi khi thêm bản ghi
		c.handleError(ctx, err)
		return
	}
	response.OkWithData(ctx, result)
}

func (c *WaterController) UpdateWaterRecords(ctx *gin.Context) {
	var requestParams request.Waterrequest

	// Kiểm tra và validate tham số yêu cầu
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		c.handleError(ctx, err)
		return
	}
	result, err := services.WaterRecord.UpdateWaterRecordSevice(&requestParams)
	if err != nil {
		// Bắt lỗi khi cập nhật bản ghi
		c.handleError(ctx, err)
		return
	}
	response.OkWithData(ctx, result)
}

func (c *WaterController) DeleteWaterRecords(ctx *gin.Context) {
	var requestParams request.Waterrequest

	// 🔹 Kiểm tra và validate tham số yêu cầu từ body
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		c.handleError(ctx, err)
		return
	}

	// 🔹 Kiểm tra RecordID có tồn tại không
	if requestParams.RecordID == "" {
		c.handleError(ctx, errors.New("thiếu RecordID, không thể xóa"))
		return
	}

	// 🔹 Gọi service để xóa
	err := services.WaterRecord.DeleteWaterRecordSevice(requestParams.RecordID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	// 🔹 Trả về response thành công
	response.OkWithData(ctx, nil)
}

func (c *WaterController) SearchWaterRecords(ctx *gin.Context) {
	var requestParams request.Waterrequest

	// Kiểm tra và validate tham số yêu cầu
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		c.handleError(ctx, err)
		return
	}
	result, err := services.WaterRecord.SearchWaterRecordSevice(&requestParams)
	if err != nil {
		// Bắt lỗi khi tìm kiếm bản ghi
		c.handleError(ctx, err)
		return
	}
	response.OkWithData(ctx, result)
}

func (c *WaterController) handleError(ctx *gin.Context, err error) {
	// Kiểm tra loại lỗi và xác định mã trạng thái HTTP phù hợp
	switch err.Error() {
	case "không tìm thấy bản ghi nào":
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
	case "dữ liệu không hợp lệ", "thiếu thông tin quan trọng":
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
	default:
		// Trả về lỗi thật để dễ debug
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
}
