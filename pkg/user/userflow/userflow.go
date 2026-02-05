package userflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/models"
	"github.com/kiranpawar037/task-managements-service/pkg/user/userflow/config"
	"gorm.io/gorm"
)

func GetUsersTasks(c *gin.Context, db *gorm.DB) {
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

	var tasks []models.Task
	if err := db.Where("user_id = ?", userobj.ID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func UpdateMyTaskStatus(c *gin.Context, db *gorm.DB) {
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

	taskID := c.Param("id")

	var req config.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	validStatus := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatus[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status. Allowed: pending, in_progress, completed",
		})
		return
	}

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", taskID, userobj.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Status = req.Status
	if err := db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task status updated successfully",
		"task": gin.H{
			"id":     task.ID,
			"status": task.Status,
		},
	})
}
