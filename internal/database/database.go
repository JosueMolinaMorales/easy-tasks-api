package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

const insertUserQuery = `
INSERT INTO Users(firstname, lastname, username, id, password, email) VALUES($1, $2, $3, $4, $5, $6);
`

// This function will make a connection to the database only once.
func init() {
	connStr := "postgres://postgres:password@localhost:5432/easy-tasks?sslmode=disable" // TODO: Make this an env var
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	log.Printf("[INFO] Connection to postgres database established") // TODO: Figure out logging
}

func InsertNewUser(newUser *types.RegisterUser) (string, error) {
	stmt, err := db.Prepare(insertUserQuery)
	if err != nil {
		return "", err
	}

	// Gen new UUID
	id := uuid.New().String()
	_, err = stmt.Exec(newUser.FirstName, newUser.LastName, strings.ToLower(newUser.Username), id, newUser.Password, strings.ToLower(newUser.Email))
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetUserByEmail(email string) (*types.User, error) {
	stmt, err := db.Prepare(`SELECT id, firstname, lastname, email, username, password FROM users WHERE email=$1`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Query(strings.ToLower(email))
	if err != nil {
		return nil, err
	}

	defer res.Close()

	user := &types.User{}
	nextVal := res.Next()
	if !nextVal {
		return nil, nil
	}

	err = res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(username string) (*types.User, error) {
	stmt, err := db.Prepare(`SELECT id, firstname, lastname, email, username, password FROM users WHERE username=$1`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Query(strings.ToLower(username))
	if err != nil {
		return nil, err
	}

	defer res.Close()
	user := &types.User{}
	nextVal := res.Next()
	if !nextVal {
		return nil, nil
	}

	err = res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(id string) (*types.User, error) {
	stmt, err := db.Prepare(`SELECT id, firstname, lastname, email, username, password FROM users WHERE id=$1`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	defer res.Close()
	user := &types.User{}
	nextVal := res.Next()
	if !nextVal {
		return nil, nil
	}

	err = res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateTask(newTask *types.Task) error {
	stmt, err := db.Prepare(`INSERT INTO Tasks(id, title, description, due_date, priority, status, created_at, updated_at, author) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newTask.ID, newTask.Title, newTask.Description, newTask.DueDate, newTask.Priority, newTask.Status, newTask.CreatedAt, newTask.UpdatedAt, newTask.Author)
	if err != nil {
		return err
	}
	return nil
}

func GetTasks(userId string) ([]*types.Task, error) {
	stmt, err := db.Prepare(`SELECT id, title, description, due_date, priority, status, created_at, updated_at, author FROM tasks WHERE author = $1`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	tasks := []*types.Task{}
	for rows.Next() {
		t := &types.Task{}
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.DueDate, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.Author)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
