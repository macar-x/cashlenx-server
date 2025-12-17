package auth_controller

import (
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/user_service"
	"github.com/macar-x/cashlenx-server/util"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login requests
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.UserLoginRequest
	if err := util.ParseJSONRequest(r, &loginRequest); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	// Validate required fields
	if loginRequest.Username == "" || loginRequest.Password == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewValidationError("username and password are required"))
		return
	}

	// Get user by username
	user := user_service.GetUserByUsername(loginRequest.Username)
	if user.Id.IsZero() {
		util.ComposeJSONResponse(w, http.StatusUnauthorized, errors.NewUnauthorizedError("invalid username or password"))
		return
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusUnauthorized, errors.NewUnauthorizedError("invalid username or password"))
		return
	}

	// TODO: Generate JWT token

	// Return user info (without password hash)
	util.ComposeJSONResponse(w, http.StatusOK, user)
}

// Register handles user registration requests
func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest model.UserRegistrationRequest
	if err := util.ParseJSONRequest(r, &registerRequest); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	// Validate required fields
	if registerRequest.Username == "" || registerRequest.Password == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewValidationError("username and password are required"))
		return
	}

	// Check if user already exists
	existingUser := user_service.GetUserByUsername(registerRequest.Username)
	if !existingUser.Id.IsZero() {
		util.ComposeJSONResponse(w, http.StatusConflict, errors.NewAlreadyExistsError("username already exists"))
		return
	}

	// Create user DTO
	userDTO := model.UserDTO{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
	}

	// Create user
	userId, err := user_service.CreateService(userDTO)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Get the created user
	createdUser := user_service.GetUserByObjectId(userId)
	if createdUser.Id.IsZero() {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError("failed to create user", nil))
		return
	}

	util.ComposeJSONResponse(w, http.StatusCreated, createdUser)
}