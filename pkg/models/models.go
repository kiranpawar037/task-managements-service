package models

import "gorm.io/gorm"

type PendingUser struct {
	gorm.Model
	Email    string `json:"email" gorm:"size:255;uniqueIndex"`
	Password string `json:"password" gorm:"size:255"`
	Role     string `json:"role" gorm:"size:20"`
}

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"size:255;uniqueIndex"`
	Password string `json:"password" gorm:"size:255"`
	Role     string `json:"role" gorm:"size:20"`
}

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // pending | in_progress | completed
	UserID      *uint  `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID"`
}
