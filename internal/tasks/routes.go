package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BuildTaskRoutes(r *gin.Engine) {
	tasksGroup := r.Group("/tasks")
	// create tasks
	tasksGroup.POST("", createTasksHandler)
}

// TODO: Add authorization to this route
func createTasksHandler(ctx *gin.Context) {
	var newTask types.CreateTask
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, fmt.Sprintf("Failed to convert JSON: %s", err.Error())))
		return
	}

	// Transform the due date
	dueDate := time.Unix(int64(newTask.DueDate), 0)
	// Validate the due date
	now := time.Now()
	// Due Date must be after the current time
	if !now.Before(dueDate) {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, "Due date must be in the future"))
		return
	}
	// Get the created_at date
	createdAt := time.Now()

	// Create the id
	id := uuid.New().String()

	// Validate the status & priority
	if newTask.Priority != types.PriorityHigh && newTask.Priority != types.PriorityLow && newTask.Priority != types.PriorityMedium {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, "Priority must be: low, medium or high"))
		return
	}
	if newTask.Status != types.StatusComplete && newTask.Status != types.StatusInProgress && newTask.Status != types.StatusPending {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusBadRequest, "Status must be: pending, in progress or complete"))
		return
	}

	task := &types.Task{
		ID:          id,
		Title:       newTask.Title,
		Description: newTask.Description,
		DueDate:     dueDate,
		Priority:    newTask.Priority,
		Status:      newTask.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	err := database.CreateTask(task)
	if err != nil {
		errors.HandleError(ctx, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to create task: %s", err.Error())))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task Created",
	})
}
