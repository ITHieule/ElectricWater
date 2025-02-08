package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type FactoriesService struct {
	*BaseService
}

var Factor = &FactoriesService{}

func (s *FactoriesService) FactoriesSevice() ([]types.Factories, error) {
	var Fac []types.Factories

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
		SELECT * FROM Factories

	`
	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&Fac).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Fac, nil
}
