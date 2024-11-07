package server

import (
	"net/http"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/auth"
	_ "github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := BuildRouter()
	r.Run("0.0.0.0:3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func BuildRouter() *gin.Engine {
	r := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	welcomeRoute(r)
	auth.BuildAuthRoutes(r)
	tasks.BuildTaskRoutes(r)

	return r
}

func welcomeRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to EasyTasks API v0.0.1!",
		})
	})
}
