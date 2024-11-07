package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/config"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/database"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/types"
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type authResponse struct {
	Token       string          `json:"token"`
	User        *types.AuthUser `json:"user"`
	GravatarURL string          `json:"gravatar_url"`
}

func registerUser(newUser *types.RegisterUser) (*authResponse, *errors.RequestError) {
	// Hash and Salt password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10) // TODO: Add a error handling middleware?
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	newUser.Password = string(hashed)

	// Validate that username & email is unique
	userUsername, err := database.GetUserByUsername(newUser.Username)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil, errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	if userUsername != nil {
		return nil, errors.NewRequestError(http.StatusBadRequest, "Username is taken")
	}

	userEmail, err := database.GetUserByEmail(newUser.Email)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil, errors.NewRequestError(http.StatusInternalServerError, "Internal Service Error")
	}
	if userEmail != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, "Email is taken")
	}

	// Insert into db
	id, err := database.InsertNewUser(newUser)
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Create JWT Token to return back along with the id/user obj
	token, err := config.NewClaims(id).SignToken()
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to create jwt token: %s", err.Error()))
	}
	// TODO: Send email notification

	return &authResponse{
		Token: token,
		User: &types.AuthUser{
			ID:        id,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Username:  newUser.Username,
			Email:     newUser.Email,
		},
		GravatarURL: utils.NewGravatarFromEmail(newUser.Email).GetURL(),
	}, nil
}

func login(loginInfo *loginInfo) (*authResponse, error) {
	if loginInfo.Email == "" && loginInfo.Username == "" {
		return nil, errors.NewRequestError(http.StatusBadRequest, "Username or Email required")
	}

	// Get user by email if provided
	var user *types.User
	var err error
	if loginInfo.Email != "" {
		user, err = database.GetUserByEmail(loginInfo.Email)
	} else {
		user, err = database.GetUserByUsername(loginInfo.Username)
	}

	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %s", err.Error()))
	}
	if user == nil {
		return nil, errors.NewRequestError(http.StatusBadRequest, "Username/Email or password is incorrect")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		return nil, errors.NewRequestError(http.StatusBadRequest, "Username/Email or password is incorrect")
	}

	// Generate token
	token, err := config.NewClaims(user.ID).SignToken()
	if err != nil {
		return nil, errors.NewRequestError(http.StatusInternalServerError, fmt.Sprintf("Failed to create jwt: %s", err.Error()))
	}

	// TODO: Update Last Login Time

	return &authResponse{
		Token: token,
		User: &types.AuthUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Username:  user.Username,
		},
		GravatarURL: utils.NewGravatarFromEmail(user.Email).GetURL(),
	}, nil
}
