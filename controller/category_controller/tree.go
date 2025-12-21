package category_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

func Tree(c *gin.Context) {
	// Extract user ID from context
	userId, exists := c.Get("user_id")
	if !exists {
		util.ResponseError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userStrId, ok := userId.(string)
	if !ok {
		util.ResponseError(c, http.StatusUnauthorized, "invalid user ID format")
		return
	}

	// Get category type from query parameter
	categoryType := c.Query("type")

	// Validate category type if provided
	if categoryType != "" && categoryType != "income" && categoryType != "expense" {
		util.ResponseError(c, http.StatusBadRequest, "category type must be 'income' or 'expense'")
		return
	}

	// Get category tree with user ID and type filter
	tree, err := category_service.TreeService(userStrId, categoryType)
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseSuccess(c, tree)
}