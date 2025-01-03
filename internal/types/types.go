package types

type User struct {
	ID        string
	FirstName string
	LastName  string
	Password  string
	Username  string
	Email     string
}

type AuthUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type RegisterUser struct {
	FirstName string `json:"firstname" binding:"min=3,max=100,required"`
	LastName  string `json:"lastname" binding:"min=3,max=100,required"`
	Password  string `json:"password" binding:"min=3,max=100,required"`
	Username  string `json:"username" binding:"min=3,max=100,required"`
	Email     string `json:"email" binding:"email,required"`
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in progress"
	StatusComplete   Status = "complete"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          string   `json:"id"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	DueDate     int      `json:"due_date"`
	Priority    Priority `json:"priority"`
	Status      Status   `json:"status"`
	CreatedAt   int      `json:"created_at"`
	UpdatedAt   int      `json:"updated_at"`
}

type CreateTask struct {
	Title       string   `json:"title" binding:"min=3,max=100,required"`
	Description string   `json:"description" binding:"min=3,max=1000,required"`
	DueDate     int      `json:"due_date" binding:"required"`
	Priority    Priority `json:"priority" binding:"required"`
	Status      Status   `json:"status" binding:"required"`
}
