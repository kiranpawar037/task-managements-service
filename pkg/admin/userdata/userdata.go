package userdata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/models"
	"gorm.io/gorm"
)

func GetAllUserByAdmin(c *gin.Context, db *gorm.DB) {
	// Get user from JWT
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userobj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}

	// Check admin role
	if userobj.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin only"})
		return
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	response := make([]gin.H, 0)
	for _, u := range users {
		response = append(response, gin.H{
			"id":         u.ID,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": response,
	})
}
