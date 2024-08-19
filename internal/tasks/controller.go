package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/google/uuid"
)

func createTask(author string, newTask *types.CreateTask) (string, *errors.RequestError) {
	// Validate the user exists
	user, err := database.GetUserByID(author)
	if err != nil {
		return "", errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %s", err.Error()))
	}
	if user == nil {
		return "", errors.NewRequestError(http.StatusNotFound, "User not found")
	}
	// Transform the due date
	dueDate := time.Unix(int64(newTask.DueDate), 0)
	// Validate the due date
	now := time.Now()
	// Due Date must be after the current time
	if !now.Before(dueDate) {
		return "", errors.NewRequestError(http.StatusBadRequest, "Due date must be in the future")
	}
	// Get the created_at date
	createdAt := time.Now()

	// Create the id
	id := uuid.New().String()

	// Validate the status & priority
	if newTask.Priority != types.PriorityHigh && newTask.Priority != types.PriorityLow && newTask.Priority != types.PriorityMedium {
		return "", errors.NewRequestError(http.StatusBadRequest, "Priority must be: low, medium or high")
	}
	if newTask.Status != types.StatusComplete && newTask.Status != types.StatusInProgress && newTask.Status != types.StatusPending {
		return "", errors.NewRequestError(http.StatusBadRequest, "Status must be: pending, in progress or complete")
	}

	task := &types.Task{
		ID:          id,
		Author:      author,
		Title:       newTask.Title,
		Description: newTask.Description,
		DueDate:     newTask.DueDate,
		Priority:    newTask.Priority,
		Status:      newTask.Status,
		CreatedAt:   int(createdAt.Unix()),
		UpdatedAt:   int(createdAt.Unix()),
	}

	err = database.CreateTask(task)
	if err != nil {
		return "", errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to create task: %s", err.Error()))
	}

	return id, nil
}

func getTasks(author string) ([]*types.Task, *errors.RequestError) {
	// Validate the user exists
	user, err := database.GetUserByID(author)
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %s", err.Error()))
	}
	if user == nil {
		return nil, errors.NewRequestError(http.StatusNotFound, "User not found")
	}

	tasks, err := database.GetTasks(author)
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to get tasks: %s", err.Error()))
	}

	return tasks, nil
}
