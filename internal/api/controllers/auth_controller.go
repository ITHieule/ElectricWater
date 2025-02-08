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
	isAuthenticated, err := services.AuthenticateUser(buser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if isAuthenticated {
		c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai thông tin đăng nhập"})
	}
}
