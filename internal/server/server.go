package server

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/auth"
	_ "github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	auth, err := auth.New()
	if err != nil {
		// TODO: Handler error
		panic("Could not create authenticator.. Panicing")
	}
	r := BuildRouter(auth)
	r.Run("0.0.0.0:3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func BuildRouter(authenticator *auth.Authenticator) *gin.Engine {
	r := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT, GET, POST, OPTION, PATCH"},
		AllowHeaders:     []string{"access-control-allow-origin, Content-Type"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("auth-session", store))

	welcomeRoute(r)
	auth.BuildAuthRoutes(r, authenticator)
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
