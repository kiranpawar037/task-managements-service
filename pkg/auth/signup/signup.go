package signup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/auth/config"
	"github.com/kiranpawar037/task-managements-service/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context, db *gorm.DB) {
	var payload config.SignupRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signup payload"})
		return
	}

	var existing models.User
	if err := db.Where("email = ?", payload.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create pending user
	pending := models.PendingUser{
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     payload.Role,
	}

	if err := db.Create(&pending).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pending user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Signup successful. Continue to login.",
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
