package user_service

import (
	"errors"

	"github.com/macar-x/cashlenx-server/mapper/user_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	// Create the user entity - always create as normal user
	// Admin users can only be created during system initialization
	userEntity := model.UserEntity{
		Id:           primitive.NewObjectID(),
		Username:     requestBody.Username,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		Role:         model.UserRoleUser, // Always create as normal user
		CreatedAt:    util.GetCurrentTime(),
		UpdatedAt:    util.GetCurrentTime(),
	}

	// Insert the user into the database
	userId := user_mapper.INSTANCE.InsertUserByEntity(userEntity)
	if userId == "" {
		return "", errors.New("failed to create user")
	}

	return userId, nil
}
