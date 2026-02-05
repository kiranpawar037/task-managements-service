package config

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"` //in_progress, completed
}
