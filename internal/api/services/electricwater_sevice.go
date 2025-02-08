package services

import (
	"fmt"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/go-sql-driver/mysql"
)

type ElectricWaterService struct {
	*BaseService
}

var Energy = &ElectricWaterService{}

func (s *ElectricWaterService) GetEnergySevice() ([]types.EnergyRecord, error) {
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

func (s *ElectricWaterService) AddEnergySevice(requestParams *request.Energyrequest) ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// 🛠 Debug đầu vào
	fmt.Printf("🔍 Request Params: %+v\n", requestParams)

	// 🛠 Tạo RecordID theo format "W" + FactoryID + Năm + Tháng
	recordID := fmt.Sprintf("E%s%d%02d", requestParams.FactoryID, requestParams.RecordYear, requestParams.RecordMonth)

	// 🛠 Load múi giờ Việt Nam
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05") // 🔹 Định dạng phù hợp với MySQL DATETIME

	// 🛠 Kết nối database
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("❌ Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🛠 Kiểm tra xem FactoryID có tồn tại trong bảng Factories không
	var factoryCount int
	checkFactoryQuery := `SELECT COUNT(*) FROM Factories WHERE FactoryID = ?`
	if err := db.Raw(checkFactoryQuery, requestParams.FactoryID).Scan(&factoryCount).Error; err != nil {
		fmt.Println("❌ Error checking FactoryID:", err)
		return nil, err
	}

	if factoryCount == 0 {
		errMsg := fmt.Sprintf("❌ FactoryID '%s' không tồn tại trong bảng Factories.", requestParams.FactoryID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Kiểm tra giá trị đầu vào
	fmt.Printf("⚡ GridElectricityMeter: %v, SolarEnergyMeter: %v\n", requestParams.GridElectricityMeter, requestParams.SolarEnergyMeter)

	// 🛠 Kiểm tra xem RecordID đã tồn tại chưa (Tránh lỗi trùng khóa)
	var existingCount int
	checkRecordQuery := `SELECT COUNT(*) FROM EnergyRecords WHERE RecordID = ?`
	if err := db.Raw(checkRecordQuery, recordID).Scan(&existingCount).Error; err != nil {
		fmt.Println("❌ Error checking existing RecordID:", err)
		return nil, err
	}

	if existingCount > 0 {
		errMsg := fmt.Sprintf("❌ RecordID '%s' đã tồn tại. Vui lòng không thêm trùng.", recordID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Thực hiện truy vấn INSERT
	query := `
        INSERT INTO EnergyRecords (
            RecordID, FactoryID, RecordYear, RecordMonth, GridElectricityMeter, 
            SolarEnergyMeter, UserID, UserDate
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	result := db.Exec(query,
		recordID,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		currentTime, // 🔹 Chuyển thành chuỗi "YYYY-MM-DD HH:MM:SS"
	)

	// 🛠 Kiểm tra lỗi truy vấn
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			errMsg := fmt.Sprintf("❌ Lỗi trùng khóa chính: RecordID '%s' đã tồn tại.", recordID)
			fmt.Println(errMsg)
			return nil, fmt.Errorf(errMsg)
		}
		fmt.Println("❌ Query execution error:", result.Error)
		return nil, result.Error
	}

	// 🛠 Kiểm tra xem có dữ liệu được thêm không
	if result.RowsAffected == 0 {
		errMsg := "❌ Không có bản ghi nào được thêm vào."
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Lấy dữ liệu vừa thêm từ database
	querySelect := `SELECT * FROM EnergyRecords WHERE RecordID = ?`
	if err := db.Raw(querySelect, recordID).Scan(&Ener).Error; err != nil {
		fmt.Println("❌ Error fetching inserted record:", err)
		return nil, err
	}

	fmt.Println("✅ Thêm dữ liệu thành công vào EnergyRecords!")
	return Ener, nil
}

func (s *ElectricWaterService) UpdateEnergySevice(requestParams *request.Energyrequest) ([]types.EnergyRecord, error) {
	var Ener []types.EnergyRecord

	// 🛠 Kiểm tra xem RecordID có tồn tại không
	var recordCount int
	db, err := database.ElectricWaterDBConnection()
	if err != nil {
		fmt.Println("❌ Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	checkRecordQuery := `SELECT COUNT(*) FROM EnergyRecords WHERE RecordID = ?`
	if err := db.Raw(checkRecordQuery, requestParams.RecordID).Scan(&recordCount).Error; err != nil {
		fmt.Println("❌ Error checking RecordID:", err)
		return nil, err
	}

	if recordCount == 0 {
		errMsg := fmt.Sprintf("❌ RecordID '%s' không tồn tại. Không thể cập nhật!", requestParams.RecordID)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Load múi giờ Việt Nam
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println("❌ Error loading time location:", err)
		return nil, err
	}
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05") // 🔹 Định dạng cho MySQL DATETIME

	// 🛠 Cập nhật dữ liệu
	query := `
		UPDATE EnergyRecords SET 
			FactoryID = ?, RecordYear = ?, RecordMonth = ?, 
			GridElectricityMeter = ?, SolarEnergyMeter = ?, 
			UserID = ?, UserDate = ?
		WHERE RecordID = ?
	`

	result := db.Exec(query,
		requestParams.FactoryID,
		requestParams.RecordYear,
		requestParams.RecordMonth,
		requestParams.GridElectricityMeter,
		requestParams.SolarEnergyMeter,
		requestParams.UserID,
		currentTime,
		requestParams.RecordID,
	)

	// 🛠 Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("❌ Query execution error:", result.Error)
		return nil, result.Error
	}

	// 🛠 Kiểm tra xem có dữ liệu nào được cập nhật không
	if result.RowsAffected == 0 {
		errMsg := "❌ Không có bản ghi nào được cập nhật. Có thể dữ liệu không thay đổi!"
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 🛠 Lấy dữ liệu vừa cập nhật từ database
	querySelect := `SELECT * FROM EnergyRecords WHERE RecordID = ?`
	if err := db.Raw(querySelect, requestParams.RecordID).Scan(&Ener).Error; err != nil {
		fmt.Println("❌ Error fetching updated record:", err)
		return nil, err
	}

	fmt.Println("✅ Cập nhật dữ liệu thành công vào EnergyRecords!")
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
