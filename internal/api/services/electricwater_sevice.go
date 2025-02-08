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
	query := `
		SELECT * FROM ElectricWaterDB.EnergyRecords

	`
	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&Ener).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
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

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `
		 INSERT INTO ElectricWaterDB.EnergyRecords (
          FactoryID,RecordYear,RecordMonth,GridElectricityMeter,SolarEnergyMeter,UserID,UserDate
        ) VALUES (?, ?,?,?,?,?,?)

	`

	// Truyền tham số vào câu truy vấn
	err = db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		requestParams.UserDate,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
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

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `
    UPDATE ElectricWaterDB.EnergyRecords
    SET 
        FactoryID = ?, 
        RecordYear = ?, 
        RecordMonth = ?, 
        GridElectricityMeter = ?, 
        SolarEnergyMeter = ?, 
        UserID = ?, 
        UserDate = ?
    WHERE RecordID = ?
`

	err = db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		requestParams.UserDate,
		requestParams.RecordID, // Thêm ID của bản ghi cần cập nhật
	).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Ener, nil
}
func (s *ElectricWaterService) DeleteEnergySevice(RecordID int) error {

	// Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM ElectricWaterDB.EnergyRecords WHERE RecordID = ?", RecordID)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No EnergyRecords found with ID:", RecordID)
		return fmt.Errorf("không tìm thấy EnergyRecords với ID %d", RecordID)
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

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán

	err = db.Raw("SELECT * FROM ElectricWaterDB.EnergyRecords WHERE FactoryID = ? OR RecordYear = ? OR RecordMonth = ? OR UserID = ?",

		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.UserID, // Đổi từ Id sang RecordID để phù hợp với bảng
	).Scan(&Ener).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Ener, nil
}
