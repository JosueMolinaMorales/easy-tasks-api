package server

import (
	"net/http"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/auth"
	_ "github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/tasks"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	welcomeRoute(r)
	auth.BuildAuthRoutes(r)
	tasks.BuildTaskRoutes(r)

	r.Run("0.0.0.0:3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func welcomeRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to EasyTasks API v0.0.1!",
		})
	})
}
