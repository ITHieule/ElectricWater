package services

import (
	"errors"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
)

// AuthenticateUser kiểm tra USERID/PWD trực tiếp trong database
func AuthenticateUser(buser request.Buser) (bool, error) {
	// Chọn database theo DBChoice
	db, err := database.SelectDB(buser.DBChoice)
	if err != nil {
		return false, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()
	// Kiểm tra user có tồn tại không
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM Busers WHERE USERID = ? AND PWD = ?", buser.UserID, buser.Password).Scan(&count).Error
	if err != nil {
		return false, errors.New("Lỗi truy vấn database")
	}

	// Nếu count > 0 nghĩa là user hợp lệ
	if count > 0 {
		return true, nil
	}

	return false, errors.New("Sai USERID hoặc PWD")
}
