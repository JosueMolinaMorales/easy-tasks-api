package auth

import (
	"log"
	"net/http"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginInfo struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func BuildAuthRoutes(r *gin.Engine, authenticator *Authenticator) {
	authRoutes := r.Group("/auth")
	authRoutes.POST("/users", registerHandler)
	authRoutes.GET("/user", isLoggedIn)
	authRoutes.GET("/login", LoginHandler(authenticator))
	authRoutes.GET("/callback", CallbackHandler(authenticator))
}

func isLoggedIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": "Not logged in",
		})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func registerHandler(ctx *gin.Context) {
	var newUser types.RegisterUser
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := registerUser(&newUser)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func loginHandler(ctx *gin.Context) {
	loginInfo := &loginInfo{}
	if err := ctx.ShouldBindJSON(loginInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("[ERROR] %s", err.Error())
		return
	}

	res, err := login(loginInfo)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
