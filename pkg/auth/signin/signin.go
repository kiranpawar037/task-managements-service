package signin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/auth/config"
	"github.com/kiranpawar037/task-managements-service/pkg/helper/jwthelper"
	"github.com/kiranpawar037/task-managements-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context, db *gorm.DB) {
	var payload config.LoginReq

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login payload"})
		return
	}

	var pending models.PendingUser
	err := db.Where("email = ?", payload.Email).First(&pending).Error

	if err == nil {
		// Found in pending → this is first-time login
		// Validate password
		if err := bcrypt.CompareHashAndPassword([]byte(pending.Password), []byte(payload.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		user := models.User{
			Email:    pending.Email,
			Password: pending.Password,
			Role:     pending.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		db.Delete(&pending)

		token, _ := jwthelper.GenerateJWTToken(user.Email)

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
		return
	}

	var user models.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT for existing user
	token, _ := jwthelper.GenerateJWTToken(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func CheckPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
