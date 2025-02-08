package services

import (
	"errors"
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type WaterRecordsService struct {
	*BaseService
}

var WaterRecord = &WaterRecordsService{}

// 🔹 Lấy toàn bộ danh sách WaterRecords
func (s *WaterRecordsService) WaterRecordSevice() ([]types.WaterRecords, error) {
	var records []types.WaterRecords

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực hiện truy vấn
	query := `SELECT * FROM WaterRecords`
	err = db.Raw(query).Scan(&records).Error
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn dữ liệu: %w", err)
	}

	return records, nil
}

// 🔹 Thêm một bản ghi mới vào WaterRecords
func (s *WaterRecordsService) AddWaterRecordSevice(requestParams *request.Waterrequest) ([]types.WaterRecords, error) {
	var records []types.WaterRecords

	// Kiểm tra tham số đầu vào
	if requestParams.FactoryID == "" || requestParams.RecordYear == 0 || requestParams.RecordMonth == 0 {
		return nil, errors.New("dữ liệu không hợp lệ, thiếu FactoryID, RecordYear hoặc RecordMonth")
	}

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Chạy câu lệnh INSERT
	query := `
		INSERT INTO WaterRecords (
			FactoryID, RecordYear, RecordMonth, TapWaterMeter, RecycledWaterMeter, UserID, UserDate
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result := db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.TapWaterMeter,
		requestParams.RecycledWaterMeter,
		requestParams.UserID,
		requestParams.UserDate,
	)

	if result.Error != nil {
		return nil, fmt.Errorf("lỗi khi thêm bản ghi: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("không có bản ghi nào được thêm")
	}

	return records, nil
}

// 🔹 Cập nhật thông tin bản ghi WaterRecords
func (s *WaterRecordsService) UpdateWaterRecordSevice(requestParams *request.Waterrequest) ([]types.WaterRecords, error) {
	var records []types.WaterRecords

	// Kiểm tra tham số đầu vào
	if requestParams.RecordID == "" {
		return nil, errors.New("thiếu RecordID, không thể cập nhật")
	}

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Chạy câu lệnh UPDATE
	query := `
		UPDATE WaterRecords
		SET FactoryID = ?, RecordYear = ?, RecordMonth = ?, TapWaterMeter = ?, RecycledWaterMeter = ?, UserID = ?, UserDate = ?
		WHERE RecordID = ?
	`
	result := db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.TapWaterMeter,
		requestParams.RecycledWaterMeter,
		requestParams.UserID,
		requestParams.UserDate,
		requestParams.RecordID,
	)

	if result.Error != nil {
		return nil, fmt.Errorf("lỗi khi cập nhật bản ghi: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("không tìm thấy bản ghi để cập nhật")
	}

	return records, nil
}

// 🔹 Xóa một bản ghi theo RecordID
func (s *WaterRecordsService) DeleteWaterRecordSevice(RecordID string) error {
	if RecordID == "" {
		return errors.New("thiếu RecordID, không thể xóa")
	}

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		return fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Chạy câu lệnh DELETE
	result := db.Exec("DELETE FROM WaterRecords WHERE RecordID = ?", RecordID)

	if result.Error != nil {
		return fmt.Errorf("lỗi khi xóa bản ghi: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("không tìm thấy bản ghi với ID %d để xóa", RecordID)
	}

	return nil
}

// 🔹 Tìm kiếm bản ghi theo điều kiện
func (s *WaterRecordsService) SearchWaterRecordSevice(requestParams *request.Waterrequest) ([]types.WaterRecords, error) {
	var records []types.WaterRecords

	// Kiểm tra tham số đầu vào
	if requestParams.FactoryID == "" && requestParams.RecordYear == 0 && requestParams.RecordMonth == 0 && requestParams.UserID == "" {
		return nil, errors.New("cần ít nhất một tiêu chí tìm kiếm")
	}

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Xây dựng câu truy vấn động dựa trên các tham số đầu vào
	query := "SELECT * FROM WaterRecords WHERE 1=1"
	var queryParams []interface{}

	if requestParams.FactoryID != "" {
		query += " AND FactoryID = ?"
		queryParams = append(queryParams, requestParams.FactoryID)
	}
	if requestParams.RecordYear != 0 {
		query += " AND RecordYear = ?"
		queryParams = append(queryParams, requestParams.RecordYear)
	}
	if requestParams.RecordMonth != 0 {
		query += " AND RecordMonth = ?"
		queryParams = append(queryParams, requestParams.RecordMonth)
	}
	if requestParams.UserID != "" {
		query += " AND UserID = ?"
		queryParams = append(queryParams, requestParams.UserID)
	}

	// Thực hiện truy vấn
	err = db.Raw(query, queryParams...).Scan(&records).Error
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tìm kiếm bản ghi: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("không tìm thấy bản ghi nào")
	}

	return records, nil
}
