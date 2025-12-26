package category_service

import (
	"errors"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User-specific operations for data isolation

// QueryByIdForUser retrieves a category by ID, ensuring it belongs to the user
func QueryByIdForUser(plainId string, userId string) (model.CategoryEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CategoryEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CategoryEntity{}, errors.New("invalid user ID")
	}

	categoryEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(plainId, userObjectId)
	if categoryEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category not found or access denied")
	}
	return categoryEntity, nil
}

// QueryByNameForUser retrieves a category by name, ensuring it belongs to the user
func QueryByNameForUser(categoryName string, userId string) (model.CategoryEntity, error) {
	// Validate category name
	if err := validation.ValidateCategoryName(categoryName); err != nil {
		return model.CategoryEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CategoryEntity{}, errors.New("invalid user ID")
	}

	categoryEntity := category_mapper.INSTANCE.GetCategoryByNameAndUser(categoryName, userObjectId)
	if categoryEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category not found or access denied")
	}
	return categoryEntity, nil
}

// QueryAllForUser queries all categories for a user with optional filtering and pagination
func QueryAllForUser(userId string, categoryType string, limit int, offset int) ([]model.CategoryEntity, int64, error) {
	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, 0, errors.New("invalid user ID")
	}

	// Get total count for this user
	totalCount := category_mapper.INSTANCE.CountAllCategoriesByUser(userObjectId)

	// Get paginated results for this user
	var categories []model.CategoryEntity

	if categoryType != "" {
		// Filter by type if provided
		categoriesByType, err := category_mapper.INSTANCE.GetCategoriesByUserAndType(userObjectId, categoryType, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		categories = categoriesByType
	} else {
		// Get all categories for user
		categories = category_mapper.INSTANCE.GetAllCategoriesByUser(userObjectId, limit, offset)
	}

	return categories, totalCount, nil
}

// DeleteByIdForUser deletes a category by ID, ensuring it belongs to the user
func DeleteByIdForUser(plainId string, userId string) (model.CategoryEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CategoryEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CategoryEntity{}, errors.New("invalid user ID")
	}

	// Check if it exists and belongs to user
	existCategoryEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(plainId, userObjectId)
	if existCategoryEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category not found or access denied")
	}

	// Delete it
	deletedEntity := category_mapper.INSTANCE.DeleteCategoryByObjectIdAndUser(plainId, userObjectId)
	if deletedEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category delete failed")
	}
	return deletedEntity, nil
}

// UpdateByIdForUser updates a category record by ID, ensuring it belongs to the user
func UpdateByIdForUser(plainId, name, categoryType, remark string, parentId string, userId string) (model.CategoryEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CategoryEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CategoryEntity{}, errors.New("invalid user ID")
	}

	// Validate optional fields if provided
	if name != "" {
		if err := validation.ValidateCategoryName(name); err != nil {
			return model.CategoryEntity{}, err
		}
	}

	// Query existing record - ensure it belongs to the user
	existingEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(plainId, userObjectId)
	if existingEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category not found or access denied")
	}

	// Update fields that are provided
	if name != "" {
		existingEntity.Name = name
	}

	if categoryType != "" {
		existingEntity.Type = categoryType
	}

	if remark != "" {
		existingEntity.Remark = remark
	}

	if parentId != "" {
		parentObjectId := util.Convert2ObjectId(parentId)
		if parentObjectId != primitive.NilObjectID {
			// Verify parent category belongs to same user
			parentEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(parentId, userObjectId)
			if parentEntity.IsEmpty() {
				return model.CategoryEntity{}, errors.New("parent category not found or access denied")
			}
			existingEntity.ParentId = parentObjectId
		}
	}

	// Update modify time
	existingEntity.ModifyTime = time.Now().UTC()

	// Call mapper to update the record
	updatedEntity := category_mapper.INSTANCE.UpdateCategoryByEntityAndUser(plainId, existingEntity, userObjectId)
	if updatedEntity.IsEmpty() {
		return model.CategoryEntity{}, errors.New("failed to update category")
	}

	return updatedEntity, nil
}

// CreateForUser creates a new category for a specific user
func CreateForUser(name, categoryType, remark string, parentId string, userId string) (model.CategoryEntity, error) {
	// Validate required fields
	if err := validation.ValidateCategoryName(name); err != nil {
		return model.CategoryEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CategoryEntity{}, errors.New("invalid user ID")
	}

	// Check if category with same name already exists for this user
	existingCategory := category_mapper.INSTANCE.GetCategoryByNameAndUser(name, userObjectId)
	if !existingCategory.IsEmpty() {
		return model.CategoryEntity{}, errors.New("category with this name already exists for this user")
	}

	// Create new category entity
	newEntity := model.CategoryEntity{
		UserId:     userObjectId,
		Name:       name,
		Type:       categoryType,
		Remark:     remark,
		CreateTime: time.Now().UTC(),
		ModifyTime: time.Now().UTC(),
	}

	// Handle parent category if provided
	if parentId != "" {
		parentObjectId := util.Convert2ObjectId(parentId)
		if parentObjectId != primitive.NilObjectID {
			// Verify parent category belongs to same user
			parentEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(parentId, userObjectId)
			if parentEntity.IsEmpty() {
				return model.CategoryEntity{}, errors.New("parent category not found or access denied")
			}
			newEntity.ParentId = parentObjectId
		}
	}

	// Insert the category
	newId := category_mapper.INSTANCE.InsertCategoryByEntity(newEntity)
	if newId == "" {
		return model.CategoryEntity{}, errors.New("failed to create category")
	}

	// Retrieve and return the created category
	createdEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(newId, userObjectId)
	return createdEntity, nil
}

// GetRootCategoriesForUser retrieves root categories (no parent) for a specific user
func GetRootCategoriesForUser(userId string, categoryType string) ([]model.CategoryEntity, error) {
	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, errors.New("invalid user ID")
	}

	var categories []model.CategoryEntity
	var err error

	if categoryType != "" {
		categories, err = category_mapper.INSTANCE.GetRootCategoriesByUserAndType(userObjectId, categoryType)
	} else {
		categories, err = category_mapper.INSTANCE.GetRootCategoriesByUser(userObjectId)
	}

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetChildCategoriesForUser retrieves child categories of a parent for a specific user
func GetChildCategoriesForUser(parentId string, userId string, categoryType string) ([]model.CategoryEntity, error) {
	// Validate parent ID
	if err := validation.ValidateID(parentId); err != nil {
		return nil, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, errors.New("invalid user ID")
	}

	// Convert parent ID
	parentObjectId := util.Convert2ObjectId(parentId)
	if parentObjectId == primitive.NilObjectID {
		return nil, errors.New("invalid parent ID")
	}

	// Verify parent category belongs to user
	parentEntity := category_mapper.INSTANCE.GetCategoryByObjectIdAndUser(parentId, userObjectId)
	if parentEntity.IsEmpty() {
		return nil, errors.New("parent category not found or access denied")
	}

	var categories []model.CategoryEntity
	var err error

	if categoryType != "" {
		categories, err = category_mapper.INSTANCE.GetCategoriesByParentIdUserAndType(parentObjectId, userObjectId, categoryType)
	} else {
		categories, err = category_mapper.INSTANCE.GetCategoriesByParentIdAndUser(parentObjectId, userObjectId)
	}

	if err != nil {
		return nil, err
	}

	return categories, nil
}
