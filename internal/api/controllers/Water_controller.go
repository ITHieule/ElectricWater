package controllers

import (
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

	// Kiểm tra và validate tham số yêu cầu
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		c.handleError(ctx, err)
		return
	}
	err := services.WaterRecord.DeleteWaterRecordSevice(requestParams.RecordID)
	if err != nil {
		// Bắt lỗi khi xóa bản ghi
		c.handleError(ctx, err)
		return
	}
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

// Hàm hỗ trợ để xử lý lỗi và trả về thông báo chi tiết
func (c *WaterController) handleError(ctx *gin.Context, err error) {
	// Kiểm tra loại lỗi và xác định mã trạng thái HTTP phù hợp
	if err.Error() == "không tìm thấy bản ghi nào" {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
	} else if err.Error() == "dữ liệu không hợp lệ" || err.Error() == "thiếu thông tin quan trọng" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
	} else {
		// Nếu là lỗi khác, trả về 500
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Có lỗi xảy ra, vui lòng thử lại.")
	}
}
