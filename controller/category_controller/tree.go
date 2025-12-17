package category_controller

import (
	"net/http"
	"strconv"

	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// Tree returns categories in a tree structure with specified depth
// GET /api/category/tree?deep={deep}&type={type}
// Parameters:
//   - deep: optional, maximum depth of the tree (0 means unlimited, default: 0)
//   - type: optional, filter by category type ("income" or "expense")
func Tree(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	deepStr := r.URL.Query().Get("deep")
	categoryType := r.URL.Query().Get("type")

	// Default depth to 0 (unlimited)
	deep := 0
	if deepStr != "" {
		if d, err := strconv.Atoi(deepStr); err == nil {
			deep = d
		}
	}

	// Validate category type if provided
	if categoryType != "" && categoryType != "income" && categoryType != "expense" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid type parameter",
			"message": "Category type must be either 'income' or 'expense'",
		})
		return
	}

	// Call service to get category tree
	categoryTree, err := category_service.TreeService(deep, categoryType)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return the tree structure
	util.ComposeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": categoryTree,
		"deep": deep,
		"type": categoryType,
	})
}
