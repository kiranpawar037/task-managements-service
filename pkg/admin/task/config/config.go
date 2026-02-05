package config

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type AssignTaskRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
