package category_service

import (
	"github.com/macar-x/cashlenx-server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TreeService builds a category tree with specified depth
// Parameters:
//   - deep: maximum depth of the tree (0 means unlimited)
//   - categoryType: optional filter for category type ("income" or "expense")
// Returns:
//   - []model.CategoryTreeNode: root categories with their children up to the specified depth
func TreeService(deep int, categoryType string) ([]model.CategoryTreeNode, error) {
	// Get all categories
	allCategories, _, err := ListAllService(0, 0) // Get all categories without pagination
	if err != nil {
		return nil, err
	}

	// Filter by type if specified
	var filteredCategories []model.CategoryEntity
	if categoryType != "" {
		for _, cat := range allCategories {
			if cat.Type == categoryType {
				filteredCategories = append(filteredCategories, cat)
			}
		}
	} else {
		filteredCategories = allCategories
	}

	// Build a map for quick parent-child lookup
	categoryMap := make(map[primitive.ObjectID]model.CategoryEntity)
	for _, cat := range filteredCategories {
		categoryMap[cat.Id] = cat
	}

	// Find all root categories (categories with no parent)
	var rootCategories []model.CategoryEntity
	for _, cat := range filteredCategories {
		if cat.ParentId == primitive.NilObjectID {
			rootCategories = append(rootCategories, cat)
		}
	}

	// Build the tree recursively for each root category
	var result []model.CategoryTreeNode
	for _, root := range rootCategories {
		rootNode := buildCategoryTree(root, categoryMap, deep, 0)
		result = append(result, rootNode)
	}

	return result, nil
}

// buildCategoryTree recursively builds the category tree up to the specified depth
func buildCategoryTree(current model.CategoryEntity, categoryMap map[primitive.ObjectID]model.CategoryEntity, maxDepth int, currentDepth int) model.CategoryTreeNode {
	// Create a new node for the current category
	node := model.NewCategoryTreeNode(current)

	// If we've reached max depth or maxDepth is 0 (unlimited), don't add children
	if maxDepth > 0 && currentDepth >= maxDepth {
		return node
	}

	// Find all children of the current category
	for _, cat := range categoryMap {
		if cat.ParentId == current.Id {
			// Recursively build children with increased depth
			childNode := buildCategoryTree(cat, categoryMap, maxDepth, currentDepth+1)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}