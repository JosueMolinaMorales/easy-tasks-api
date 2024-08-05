package tasks

import (
	"fmt"
	"net/http"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/config"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/middleware"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/gin-gonic/gin"
)

func BuildTaskRoutes(r *gin.Engine) {
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware)
	// create tasks
	tasksGroup.POST("", createTasksHandler)
}

func createTasksHandler(ctx *gin.Context) {
	var newTask types.CreateTask
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, fmt.Sprintf("Failed to convert JSON: %s", err.Error())))
		return
	}
	// Get the userID of the author
	author, err := config.ExtractIDFromToken(ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	id, err := createTask(author, &newTask)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}
