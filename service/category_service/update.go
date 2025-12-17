package category_service

import (
	"errors"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
)

// UpdateService updates a category by ID
func UpdateService(plainId, parentPlainId, categoryName, categoryType string) error {
	if plainId == "" {
		return errors.New("id cannot be empty")
	}

	// Query existing category
	existingCategory := category_mapper.INSTANCE.GetCategoryByObjectId(plainId)
	if existingCategory.IsEmpty() {
		return errors.New("category not found")
	}

	// Check if category has children before allowing type change
	categoryChildren := category_mapper.INSTANCE.GetCategoryByParentId(plainId)

	// Update parent ID if provided
	if parentPlainId != "" {
		parentCategory := category_mapper.INSTANCE.GetCategoryByObjectId(parentPlainId)
		if parentCategory.IsEmpty() {
			return errors.New("parent category does not exist")
		}
		// Prevent circular reference
		if parentPlainId == plainId {
			return errors.New("category cannot be its own parent")
		}
		// Ensure parent and child categories have the same type
		// Use existing type or new type if being updated
		childType := existingCategory.Type
		if categoryType != "" {
			childType = categoryType
		}
		if parentCategory.Type != childType {
			return errors.New("parent and child categories must have the same type")
		}
		existingCategory.ParentId = parentCategory.Id
	}

	// Update category name if provided
	if categoryName != "" {
		existingCategory.Name = categoryName
	}

	// Update category type if provided
	if categoryType != "" {
		// Validate category type
		if categoryType != "income" && categoryType != "expense" {
			return errors.New("category type must be either 'income' or 'expense'")
		}
		// Prevent type change if category has children
		if len(categoryChildren) > 0 && existingCategory.Type != categoryType {
			return errors.New("cannot change type of a category with children")
		}
		existingCategory.Type = categoryType
	}

	// Call mapper to update the record
	updatedEntity := category_mapper.INSTANCE.UpdateCategoryByEntity(plainId, existingCategory)
	if updatedEntity.IsEmpty() {
		return errors.New("failed to update category")
	}

	return nil
}
