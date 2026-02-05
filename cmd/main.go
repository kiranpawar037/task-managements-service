package main

import (
	"fmt"
	"os"

	"github.com/kiranpawar037/task-managements-service/pkg/admin/worker"
	"github.com/kiranpawar037/task-managements-service/pkg/database"
	"github.com/kiranpawar037/task-managements-service/routes/useradmin"
)

func main() {
	dbConn, err := database.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)

		return
	}
	worker.StartTaskAutoCompleteWorker(dbConn)
	var serviceName string

	// Check if the SERVICE_NAME environment variable is set
	if envServiceName := os.Getenv("SERVICE_NAME"); envServiceName != "" {
		serviceName = envServiceName
	} else {
		// Check if a command-line argument is provided
		if len(os.Args) < 2 {
			fmt.Println("Service name not provided")
			return
		}
		serviceName = os.Args[1]
	}

	switch serviceName {
	case "user-admin":
		useradmin.UserAdmin(dbConn)
	}
}

//done
