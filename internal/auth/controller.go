package auth

import (
	"log"
	"net/http"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(newUser *types.RegisterUser) (string, *errors.RequestError) {
	// Hash and Salt password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10) // TODO: Add a error handling middleware?
	if err != nil {
		return "", errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	newUser.Password = string(hashed)

	// Validate that username & email is unique
	userUsername, err := database.GetUserByUsername(newUser.Username)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return "", errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	if userUsername != nil {
		return "", errors.NewRequestError(http.StatusBadRequest, "Username is taken")
	}

	userEmail, err := database.GetUserByEmail(newUser.Email)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return "", errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	if userEmail != nil {
		return "", errors.NewRequestError(http.StatusInternalServerError, "Email is taken")
	}

	// Insert into db
	id, err := database.InsertNewUser(newUser)
	if err != nil {
		return "", errors.NewRequestError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Create JWT Token to return back along with the id/user obj

	// TODO: Send email notification

	return id, nil
}
