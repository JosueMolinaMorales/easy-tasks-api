package middleware

import (
	"net/http"
	"strings"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/config"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	// Get the headers
	headers := ctx.Request.Header
	authHeader := headers.Get("Authorization")
	// Check if auth header exists
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header required",
		})
		return
	}

	// Check if header is formatted correctly
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header formatted incorrectly",
		})
		return
	}

	if strings.ToLower(parts[0]) != "bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header formatted incorrectly",
		})
		return
	}

	token := parts[1]

	claims, err := config.VerifyToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token is invalid",
		})
		return
	}

	// Store claims
	ctx.Set("Token", claims)
}
