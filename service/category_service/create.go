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

func CreateService(parentPlainId, categoryName string) (string, error) {
	// Validate category name
	if err := validation.ValidateCategoryName(categoryName); err != nil {
		return "", err
	}

	// Validate parent ID if provided
	if parentPlainId != "" {
		if err := validation.ValidateID(parentPlainId); err != nil {
			return "", err
		}
	}

	categoryEntity := model.CategoryEntity{
		ParentId: primitive.NilObjectID,
		Name:     categoryName,
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
