package category_service

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TreeService builds a category tree with specified depth
// Parameters:
//   - deep: maximum depth of the tree (0 means unlimited)
//   - categoryType: optional filter for category type ("income" or "expense")
//
// Returns:
//   - []model.CategoryTreeNode: root categories with their children up to the specified depth
func TreeService(deep int, categoryType string) ([]model.CategoryTreeNode, error) {
	// Get all categories
	allCategories, _, err := ListAllService("", categoryType, 0, 0) // Get all categories without pagination
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

// GetCategoryTreeByUser builds category tree for a specific user
func GetCategoryTreeByUser(userId, categoryType string) ([]model.CategoryTree, error) {
	// Validate user ID
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Validate category type
	if categoryType != "income" && categoryType != "expense" && categoryType != "" {
		return nil, fmt.Errorf("category type must be 'income', 'expense', or empty")
	}

	// Get root categories with user ID and type filter
	var rootCategories []model.CategoryEntity
	var err error

	if categoryType == "" {
		rootCategories, err = category_mapper.INSTANCE.GetRootCategoriesByUser(userObjectId)
	} else {
		rootCategories, err = category_mapper.INSTANCE.GetRootCategoriesByUserAndType(userObjectId, categoryType)
	}

	if err != nil {
		return nil, err
	}

	// Build category tree recursively
	var categoryTreeList []model.CategoryTree
	for _, root := range rootCategories {
		categoryTree := buildUserCategoryTree(root, userObjectId, categoryType)
		categoryTreeList = append(categoryTreeList, categoryTree)
	}

	return categoryTreeList, nil
}

func buildUserCategoryTree(parent model.CategoryEntity, userId primitive.ObjectID, categoryType string) model.CategoryTree {
	// Convert entity to tree node
	categoryTree := model.CategoryTree{
		Id:       parent.Id.Hex(),
		ParentId: parent.ParentId.Hex(),
		Name:     parent.Name,
		Type:     parent.Type,
		Children: []model.CategoryTree{},
	}

	// Get children with user ID and type filter
	var children []model.CategoryEntity
	var err error

	if categoryType == "" {
		children, err = category_mapper.INSTANCE.GetCategoriesByParentIdAndUser(parent.Id, userId)
	} else {
		children, err = category_mapper.INSTANCE.GetCategoriesByParentIdUserAndType(parent.Id, userId, categoryType)
	}

	if err != nil {
		fmt.Printf("Error getting children for category %s: %v", parent.Id.Hex(), err)
		return categoryTree
	}

	// Recursively build children nodes
	for _, child := range children {
		childTree := buildUserCategoryTree(child, userId, categoryType)
		categoryTree.Children = append(categoryTree.Children, childTree)
	}

	return categoryTree
}
