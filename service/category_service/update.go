package category_service

import (
	"errors"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
)

// UpdateService updates a category by ID
func UpdateService(plainId, parentPlainId, categoryName string) error {
	if plainId == "" {
		return errors.New("id cannot be empty")
	}

	// Query existing category
	existingCategory := category_mapper.INSTANCE.GetCategoryByObjectId(plainId)
	if existingCategory.IsEmpty() {
		return errors.New("category not found")
	}

	// Update fields that are provided
	if parentPlainId != "" {
		parentCategory := category_mapper.INSTANCE.GetCategoryByObjectId(parentPlainId)
		if parentCategory.IsEmpty() {
			return errors.New("parent category does not exist")
		}
		// Prevent circular reference
		if parentPlainId == plainId {
			return errors.New("category cannot be its own parent")
		}
		existingCategory.ParentId = parentCategory.Id
	}

	if categoryName != "" {
		existingCategory.Name = categoryName
	}

	// Call mapper to update the record
	updatedEntity := category_mapper.INSTANCE.UpdateCategoryByEntity(plainId, existingCategory)
	if updatedEntity.IsEmpty() {
		return errors.New("failed to update category")
	}

	return nil
}
