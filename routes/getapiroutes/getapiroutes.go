package getapiroutes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kiranpawar037/task-managements-service/pkg/helper/cors"
)

func GetApiRoutes() (*gin.RouterGroup, *gin.Engine) {
	// version
	getCurrentApiVersion := os.Getenv("API_VERSION")

	if getCurrentApiVersion == "" {
		getCurrentApiVersion = "v1"
		log.Printf("Defaulting to version %s", getCurrentApiVersion)
	}

	router := gin.Default()

	router.Use(cors.CORSMiddleware())

	apiV1 := router.Group(getCurrentApiVersion)

	return apiV1, router
}
