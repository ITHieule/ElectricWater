package services

import (
	"errors"
	"fmt"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/go-sql-driver/mysql"
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
		return nil, errors.New("❌ Dữ liệu không hợp lệ, thiếu FactoryID, RecordYear hoặc RecordMonth")
	}

	// 🛠 Tạo RecordID theo format "W" + FactoryID + Năm + Tháng
	recordID := fmt.Sprintf("W%s%d%02d", requestParams.FactoryID, requestParams.RecordYear, requestParams.RecordMonth)

	// 🛠 Load múi giờ Việt Nam
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05") // 🔹 Định dạng phù hợp với MySQL DATETIME

	// 🛠 Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("❌ Lỗi kết nối database:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🛠 Kiểm tra xem FactoryID có tồn tại không
	var factoryCount int
	checkFactoryQuery := `SELECT COUNT(*) FROM Factories WHERE FactoryID = ?`
	if err := db.Raw(checkFactoryQuery, requestParams.FactoryID).Scan(&factoryCount).Error; err != nil {
		fmt.Println("❌ Lỗi kiểm tra FactoryID:", err)
		return nil, err
	}
	if factoryCount == 0 {
		errMsg := fmt.Sprintf("❌ FactoryID '%s' không tồn tại trong bảng Factories.", requestParams.FactoryID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Kiểm tra xem RecordID đã tồn tại chưa
	var existingCount int
	checkRecordQuery := `SELECT COUNT(*) FROM WaterRecords WHERE RecordID = ?`
	if err := db.Raw(checkRecordQuery, recordID).Scan(&existingCount).Error; err != nil {
		fmt.Println("❌ Lỗi kiểm tra RecordID:", err)
		return nil, err
	}
	if existingCount > 0 {
		errMsg := fmt.Sprintf("❌ RecordID '%s' đã tồn tại. Vui lòng không thêm trùng.", recordID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Thực hiện INSERT vào bảng WaterRecords
	query := `
		INSERT INTO WaterRecords (
			RecordID, FactoryID, RecordYear, RecordMonth, TapWaterMeter, 
			RecycledWaterMeter, UserID, UserDate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result := db.Exec(query,
		recordID,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.TapWaterMeter,
		requestParams.RecycledWaterMeter,
		requestParams.UserID,
		currentTime,
	)

	// 🛠 Kiểm tra lỗi truy vấn
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			errMsg := fmt.Sprintf("❌ Lỗi trùng khóa chính: RecordID '%s' đã tồn tại.", recordID)
			fmt.Println(errMsg)
			return nil, fmt.Errorf(errMsg)
		}
		fmt.Println("❌ Lỗi khi thêm bản ghi:", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		errMsg := "❌ Không có bản ghi nào được thêm vào."
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Lấy dữ liệu vừa thêm từ database
	querySelect := `SELECT * FROM WaterRecords WHERE RecordID = ?`
	if err := db.Raw(querySelect, recordID).Scan(&records).Error; err != nil {
		fmt.Println("❌ Lỗi khi lấy bản ghi đã thêm:", err)
		return nil, err
	}

	fmt.Println("✅ Thêm dữ liệu thành công vào WaterRecords!")
	return records, nil
}

// 🔹 Cập nhật thông tin bản ghi WaterRecords
func (s *WaterRecordsService) UpdateWaterRecordSevice(requestParams *request.Waterrequest) ([]types.WaterRecords, error) {
	var records []types.WaterRecords

	// 🛠 Kiểm tra tham số đầu vào
	if requestParams.RecordID == "" {
		return nil, errors.New("❌ Thiếu RecordID, không thể cập nhật")
	}

	// 🛠 Load múi giờ Việt Nam
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05") // 🔹 Định dạng phù hợp với MySQL DATETIME

	// 🛠 Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("❌ Lỗi kết nối database:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🛠 Kiểm tra xem RecordID có tồn tại không
	var existingCount int
	checkRecordQuery := `SELECT COUNT(*) FROM WaterRecords WHERE RecordID = ?`
	if err := db.Raw(checkRecordQuery, requestParams.RecordID).Scan(&existingCount).Error; err != nil {
		fmt.Println("❌ Lỗi kiểm tra RecordID:", err)
		return nil, err
	}
	if existingCount == 0 {
		errMsg := fmt.Sprintf("❌ RecordID '%s' không tồn tại, không thể cập nhật.", requestParams.RecordID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Chạy câu lệnh UPDATE
	query := `
		UPDATE WaterRecords
		SET FactoryID = ?, RecordYear = ?, RecordMonth = ?, TapWaterMeter = ?, 
			RecycledWaterMeter = ?, UserID = ?, UserDate = ?
		WHERE RecordID = ?
	`
	result := db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.TapWaterMeter,
		requestParams.RecycledWaterMeter,
		requestParams.UserID,
		currentTime,
		requestParams.RecordID,
	)

	// 🛠 Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("❌ Lỗi khi cập nhật bản ghi:", result.Error)
		return nil, fmt.Errorf("lỗi khi cập nhật bản ghi: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		errMsg := "❌ Không tìm thấy bản ghi để cập nhật hoặc dữ liệu không thay đổi."
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Truy vấn lại dữ liệu vừa cập nhật
	querySelect := `SELECT * FROM WaterRecords WHERE RecordID = ?`
	if err := db.Raw(querySelect, requestParams.RecordID).Scan(&records).Error; err != nil {
		fmt.Println("❌ Lỗi khi lấy bản ghi đã cập nhật:", err)
		return nil, err
	}

	fmt.Println("✅ Cập nhật dữ liệu thành công!")
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
		return fmt.Errorf("không tìm thấy bản ghi với ID %s để xóa", RecordID)
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

	// Xây dựng câu truy vấn động
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
