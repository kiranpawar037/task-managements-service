package worker

import (
	"log"
	"time"

	"github.com/kiranpawar037/task-managements-service/pkg/config"
	"github.com/kiranpawar037/task-managements-service/pkg/models"
	"gorm.io/gorm"
)

// Channel to receive task IDs
var TaskQueue = make(chan uint, 100)

func StartTaskAutoCompleteWorker(db *gorm.DB) {
	cfg, err := config.Env()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	delay := time.Duration(cfg.Task.AutoCompleteMinutes) * time.Minute

	go func() {
		for taskID := range TaskQueue {
			go processTask(db, taskID, delay)
		}
	}()
}

func processTask(db *gorm.DB, taskID uint, delay time.Duration) {
	time.Sleep(delay)

	var task models.Task
	err := db.First(&task, taskID).Error
	if err != nil {
		// Task deleted → do nothing
		return
	}

	// Only auto-complete if still pending or in_progress
	if task.Status == "pending" || task.Status == "in_progress" {
		task.Status = "completed"
		db.Save(&task)
	}
}
