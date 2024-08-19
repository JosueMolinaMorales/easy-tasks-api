package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/JosueMolinaMorales/EasyTasksAPI/pkg/env"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

const insertUserQuery = `
INSERT INTO Users(firstname, lastname, username, id, password, email) VALUES($1, $2, $3, $4, $5, $6);
`

// This function will make a connection to the database only once.
func init() {
	var err error
	db, err = sql.Open("postgres", env.GetDBURI())
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
	_, err = stmt.Exec(newTask.ID, newTask.Title, newTask.Description, time.Unix(int64(newTask.DueDate), 0), newTask.Priority, newTask.Status, time.Unix(int64(newTask.CreatedAt), 0), time.Unix(int64(newTask.UpdatedAt), 0), newTask.Author)
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
		var dueDate time.Time
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &dueDate, &t.Priority, &t.Status, &createdAt, &updatedAt, &t.Author)
		if err != nil {
			return nil, err
		}
		t.DueDate = int(dueDate.Unix())
		t.CreatedAt = int(createdAt.Unix())
		t.UpdatedAt = int(updatedAt.Unix())
		tasks = append(tasks, t)
	}

	return tasks, nil
}
