package category_service

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateService(userId, parentPlainId, categoryName, categoryType string) (string, error) {
	// Validate user ID
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return "", errors.New("invalid user ID")
	}

	// Validate category name
	if err := validation.ValidateCategoryName(categoryName); err != nil {
		return "", err
	}

	// Validate category type
	if categoryType != "income" && categoryType != "expense" {
		return "", errors.New("category type must be either 'income' or 'expense'")
	}

	// Validate parent ID if provided
	if parentPlainId != "" {
		if err := validation.ValidateID(parentPlainId); err != nil {
			return "", err
		}

		// Check if parent category exists and belongs to the same user
		parentCategory := category_mapper.INSTANCE.GetCategoryByObjectId(parentPlainId)
		if parentCategory.IsEmpty() {
			return "", errors.New("parent category not found")
		}

		// Ensure parent category belongs to the same user
		if parentCategory.UserId != userObjectId {
			return "", errors.New("parent category does not belong to user")
		}

		// Ensure parent and child categories have the same type
		if parentCategory.Type != categoryType {
			return "", errors.New("parent and child categories must have the same type")
		}
	}

	categoryEntity := model.CategoryEntity{
		ParentId: primitive.NilObjectID,
		UserId:   userObjectId,
		Name:     categoryName,
		Type:     categoryType,
	}
	if parentPlainId != "" {
		categoryEntity.ParentId = util.Convert2ObjectId(parentPlainId)
	}

	newCategoryPlainId := category_mapper.INSTANCE.InsertCategoryByEntity(categoryEntity)
	if newCategoryPlainId == "" {
		return "", errors.New("category create failed")
	}

	newCategoryEntity := category_mapper.INSTANCE.GetCategoryByObjectId(newCategoryPlainId)
	fmt.Println("category ", 0, ": ", newCategoryEntity.ToString())
	return newCategoryPlainId, nil
}

func isCreateRequiredFiledSatisfied(categoryName string) bool {
	return categoryName != ""
}