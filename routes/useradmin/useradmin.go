package useradmin

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/admin/task"
	"github.com/kiranpawar037/task-managements-service/pkg/admin/userdata"
	"github.com/kiranpawar037/task-managements-service/pkg/auth/signin"
	"github.com/kiranpawar037/task-managements-service/pkg/auth/signup"
	"github.com/kiranpawar037/task-managements-service/pkg/middleware"
	"github.com/kiranpawar037/task-managements-service/pkg/user/userflow"
	"github.com/kiranpawar037/task-managements-service/routes/getapiroutes"
	"gorm.io/gorm"
)

func UserAdmin(db *gorm.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10001"
		log.Printf("Defaulting to port %s", port)
	}

	apiV1, router := getapiroutes.GetApiRoutes()

	// Define handlers oauth
	apiV1.GET("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "user Service Healthy")
	})

	apiV1.POST("/user/sign-up", func(c *gin.Context) {
		signup.Signup(c, db)
	})

	apiV1.POST("/user/sign-in", func(c *gin.Context) {
		signin.Login(c, db)
	})

	//admin routes
	apiV1.POST("/admin/create-task", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		task.CreateTask(c, db)
	})

	apiV1.GET("/admin/get-all-tasks", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		task.GetAllTasks(c, db)
	})

	apiV1.GET("/admin/get-task/:id", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		task.GetTaskByID(c, db)
	})

	apiV1.DELETE("/admin/delete-task/:id", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		task.DeleteTask(c, db)
	})

	apiV1.PUT("/admin/assign-task/:id", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		task.AssignTaskToUser(c, db)
	})

	apiV1.GET("/admin/get-all-users", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		userdata.GetAllUserByAdmin(c, db)
	})

	// USER APIs
	apiV1.GET("/user/get-my-tasks", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		userflow.GetUsersTasks(c, db)
	})

	apiV1.PATCH("/user/update-task-status/:id", middleware.JWTMiddlewareUser(db), func(c *gin.Context) {
		userflow.UpdateMyTaskStatus(c, db)
	})

	// Listen and serve on defined port
	log.Printf("Application started, Listening on Port %s", port)
	router.Run(":" + port)
}
