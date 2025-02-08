package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"

	"github.com/gin-gonic/gin"
)

// API đăng nhập
func Login(c *gin.Context) {
	var buser request.Buser

	// Bind JSON từ request
	if err := c.ShouldBindJSON(&buser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Gọi hàm xác thực user
	authenticatedUser, err := services.AuthenticateUser(buser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Trả về thông tin user sau khi xác thực thành công
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message":   "Đăng nhập thành công",
			"USERID":    authenticatedUser.UserID,
			"DB_CHOICE": authenticatedUser.DBChoice,
		},
	})
}
