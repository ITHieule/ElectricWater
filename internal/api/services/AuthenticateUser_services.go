package services

import (
	"errors"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

// AuthenticateUser kiểm tra USERID/PWD và trả về thông tin user nếu hợp lệ
func AuthenticateUser(buser request.Buser) (*types.BuserType, error) {
	// Chọn database theo DBChoice từ request
	db, err := database.SelectDB(buser.DBChoice)
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Kiểm tra user có tồn tại không
	var user types.BuserType
	err = db.Raw("SELECT USERID FROM Busers WHERE USERID = ? AND PWD = ?", buser.UserID, buser.Password).
		Scan(&user.UserID).Error
	if err != nil {
		return nil, errors.New("Lỗi truy vấn database")
	}

	// Nếu UserID tồn tại, gán DBChoice từ request vào struct và trả về
	if user.UserID != "" {
		user.DBChoice = buser.DBChoice // Gán DBChoice từ request vào struct trước khi trả về
		return &user, nil
	}

	return nil, errors.New("Sai tài khoản hoặc mât khẩu")
}
