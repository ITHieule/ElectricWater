package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type ElectricWaterService struct {
	*BaseService
}

var Energy = &ElectricWaterService{}

func (s *ElectricWaterService) IceLevelsSevice() ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// Kết nối databaser
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `SELECT * FROM EnergyRecords`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	if err := db.Raw(query).Scan(&Ener).Error; err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Kiểm tra kết quả trả về có dữ liệu hay không
	if len(Ener) == 0 {
		fmt.Println("No records found")
		return nil, fmt.Errorf("no records found")
	}

	return Ener, nil
}

func (s *ElectricWaterService) AddFlavorsSevice(requestParams *request.Energyrequest) ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL
	query := `
		INSERT INTO EnergyRecords (
			RecordID,FactoryID, RecordYear, RecordMonth, GridElectricityMeter, 
			SolarEnergyMeter, UserID, UserDate
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Thực hiện truy vấn với tham số
	if err := db.Exec(query, 
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		requestParams.UserDate).Error; err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Kiểm tra việc thêm dữ liệu thành công
	return Ener, nil
}

func (s *ElectricWaterService) UpdateEnergySevice(requestParams *request.Energyrequest) ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL
	query := `
		UPDATE EnergyRecords SET 
			FactoryID = ?, RecordYear = ?, RecordMonth = ?, 
			GridElectricityMeter = ?, SolarEnergyMeter = ?, 
			UserID = ?, UserDate = ?
		WHERE RecordID = ?
	`

	// Thực hiện câu lệnh cập nhật
	if err := db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		requestParams.UserDate,
		requestParams.RecordID).Error; err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Kiểm tra việc cập nhật thành công
	return Ener, nil
}

func (s *ElectricWaterService) DeleteEnergySevice(RecordID string) error {
	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM EnergyRecords WHERE RecordID = ?", RecordID)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng
	if result.RowsAffected == 0 {
		fmt.Println("No EnergyRecords found with ID:", RecordID)
		return fmt.Errorf("no EnergyRecords found with ID %d", RecordID)
	}

	fmt.Println("Deleted EnergyRecords successfully!")
	return nil
}

func (s *ElectricWaterService) SearchEnergySevice(requestParams *request.Energyrequest) ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL
	if err := db.Raw("SELECT * FROM EnergyRecords WHERE FactoryID = ? OR RecordYear = ? OR RecordMonth = ? OR UserID = ?",
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.UserID).Scan(&Ener).Error; err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Kiểm tra kết quả tìm kiếm
	if len(Ener) == 0 {
		fmt.Println("No records found matching the search criteria")
		return nil, fmt.Errorf("no records found matching the search criteria")
	}

	return Ener, nil
}
