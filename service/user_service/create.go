package user_service

import (
	"errors"

	"github.com/macar-x/cashlenx-server/mapper/user_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateService creates a new user with the provided details
func CreateService(requestBody model.UserDTO) (string, error) {
	// Check if username is already taken
	existingUser := user_mapper.INSTANCE.GetUserByUsername(requestBody.Username)
	if !existingUser.Id.IsZero() {
		return "", errors.New("username is already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create the user entity
	userEntity := model.UserEntity{
		Id:           primitive.NewObjectID(),
		Username:     requestBody.Username,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		Role:         model.UserRoleUser,
		CreatedAt:    util.GetCurrentTime(),
		UpdatedAt:    util.GetCurrentTime(),
	}

	// Set role if provided
	if requestBody.Role != "" {
		userEntity.Role = requestBody.Role
	}

	// Insert the user into the database
	userId := user_mapper.INSTANCE.InsertUserByEntity(userEntity)
	if userId == "" {
		return "", errors.New("failed to create user")
	}

	return userId, nil
}