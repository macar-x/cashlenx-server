package user_mapper

import (
	"context"
	"time"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserMongoDbMapper MongoDB implementation of UserMapper
type UserMongoDbMapper struct{}

// GetUserByObjectId retrieves a user by their ID from MongoDB
func (m UserMongoDbMapper) GetUserByObjectId(plainId string) model.UserEntity {
	// Parse the plain ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(plainId)
	if err != nil {
		util.Logger.Errorw("Invalid ObjectID", "error", err, "plainId", plainId)
		return model.UserEntity{}
	}

	// Create a filter to find the user by ID
	filter := bson.M{"_id": objectId}

	// Get the user document from the database
	var user model.UserEntity
	err = util.database.MongoDBClient.Collection(model.TableUser).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.Logger.Debugw("User not found", "userId", plainId)
			return model.UserEntity{}
		}
		util.Logger.Errorw("Failed to get user", "error", err, "userId", plainId)
		return model.UserEntity{}
	}

	return user
}

// GetUserByUsername retrieves a user by their username from MongoDB
func (m UserMongoDbMapper) GetUserByUsername(username string) model.UserEntity {
	// Create a filter to find the user by username
	filter := bson.M{"username": username}

	// Get the user document from the database
	var user model.UserEntity
	err := util.database.MongoDBClient.Collection(model.TableUser).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.Logger.Debugw("User not found", "username", username)
			return model.UserEntity{}
		}
		util.Logger.Errorw("Failed to get user", "error", err, "username", username)
		return model.UserEntity{}
	}

	return user
}

// InsertUserByEntity inserts a new user entity into MongoDB
func (m UserMongoDbMapper) InsertUserByEntity(newEntity model.UserEntity) string {
	// Set default values if not provided
	if newEntity.CreatedAt.IsZero() {
		newEntity.CreatedAt = time.Now()
	}
	if newEntity.UpdatedAt.IsZero() {
		newEntity.UpdatedAt = time.Now()
	}
	if newEntity.IsActive == false {
		newEntity.IsActive = true
	}
	if newEntity.Role == "" {
		newEntity.Role = model.UserRoleUser
	}

	// Insert the user document into the database
	result, err := util.database.MongoDBClient.Collection(model.TableUser).InsertOne(context.TODO(), newEntity)
	if err != nil {
		util.Logger.Errorw("Failed to insert user", "error", err, "username", newEntity.Username)
		return ""
	}

	// Return the inserted ID as a string
	return result.InsertedID.(primitive.ObjectID).Hex()
}

// UpdateUserByEntity updates an existing user entity in MongoDB
func (m UserMongoDbMapper) UpdateUserByEntity(plainId string, updatedEntity model.UserEntity) model.UserEntity {
	// Parse the plain ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(plainId)
	if err != nil {
		util.Logger.Errorw("Invalid ObjectID", "error", err, "plainId", plainId)
		return model.UserEntity{}
	}

	// Create a filter to find the user by ID
	filter := bson.M{"_id": objectId}

	// Set the updated timestamp
	updatedEntity.UpdatedAt = time.Now()

	// Create an update document with the updated entity
	update := bson.M{"$set": updatedEntity}

	// Update the user document in the database
	result := util.database.MongoDBClient.Collection(model.TableUser).FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	// Decode the updated user document
	var user model.UserEntity
	err = result.Decode(&user)
	if err != nil {
		util.Logger.Errorw("Failed to update user", "error", err, "userId", plainId)
		return model.UserEntity{}
	}

	return user
}

// GetAllUsers retrieves all users with pagination from MongoDB
func (m UserMongoDbMapper) GetAllUsers(limit, offset int) []model.UserEntity {
	// Create an empty filter to get all users
	filter := bson.M{}

	// Set up pagination options
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	// Get the user documents from the database
	cursor, err := util.database.MongoDBClient.Collection(model.TableUser).Find(context.TODO(), filter, opts)
	if err != nil {
		util.Logger.Errorw("Failed to get all users", "error", err)
		return []model.UserEntity{}
	}
	defer cursor.Close(context.TODO())

	// Decode the user documents
	var users []model.UserEntity
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		util.Logger.Errorw("Failed to decode users", "error", err)
		return []model.UserEntity{}
	}

	return users
}

// CountAllUsers returns the total number of users from MongoDB
func (m UserMongoDbMapper) CountAllUsers() int64 {
	// Create an empty filter to count all users
	filter := bson.M{}

	// Count the user documents in the database
	count, err := util.database.MongoDBClient.Collection(model.TableUser).CountDocuments(context.TODO(), filter)
	if err != nil {
		util.Logger.Errorw("Failed to count users", "error", err)
		return 0
	}

	return count
}

// DeleteUserByObjectId deletes a user by their ID from MongoDB
func (m UserMongoDbMapper) DeleteUserByObjectId(plainId string) model.UserEntity {
	// Parse the plain ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(plainId)
	if err != nil {
		util.Logger.Errorw("Invalid ObjectID", "error", err, "plainId", plainId)
		return model.UserEntity{}
	}

	// Create a filter to find the user by ID
	filter := bson.M{"_id": objectId}

	// Delete the user document from the database and return the deleted document
	result := util.database.MongoDBClient.Collection(model.TableUser).FindOneAndDelete(context.TODO(), filter)

	// Decode the deleted user document
	var user model.UserEntity
	err = result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.Logger.Debugw("User not found for deletion", "userId", plainId)
			return model.UserEntity{}
		}
		util.Logger.Errorw("Failed to delete user", "error", err, "userId", plainId)
		return model.UserEntity{}
	}

	return user
}

// TruncateUsers deletes all users from MongoDB
func (m UserMongoDbMapper) TruncateUsers() error {
	// Create an empty filter to delete all users
	filter := bson.M{}

	// Delete all user documents from the database
	result, err := util.database.MongoDBClient.Collection(model.TableUser).DeleteMany(context.TODO(), filter)
	if err != nil {
		util.Logger.Errorw("Failed to truncate users", "error", err)
		return err
	}

	util.Logger.Infow("Truncated users collection", "deleted_count", result.DeletedCount)
	return nil
}