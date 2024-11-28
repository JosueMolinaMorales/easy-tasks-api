package tasks

import (
	"fmt"
	"log"
	"net/http"

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
	// getting tasks
	tasksGroup.GET("", getTasksHandler)
}

func createTasksHandler(ctx *gin.Context) {
	var newTask types.CreateTask
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, fmt.Sprintf("Failed to convert JSON: %s", err.Error())))
		return
	}
	// Get the userID of the author
	profile := ctx.GetStringMap(middleware.PROFILE_KEY)
	if profile == nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusUnauthorized, "Not Logged In"))
		return
	}

	log.Printf("[DEBUG] sub : %s", profile["sub"])
	id, err := createTask(profile["sub"].(string), &newTask)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func getTasksHandler(ctx *gin.Context) {
	// Get the userID of the author
	profile := ctx.GetStringMap(middleware.PROFILE_KEY)
	if profile == nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusUnauthorized, "Not Logged In"))
		return
	}

	log.Printf("[DEBUG] sub : %s", profile["sub"])

	author := profile["sub"].(string)
	tasks, err := getTasks(author)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
