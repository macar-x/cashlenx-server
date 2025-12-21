package category_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

func Create(c *gin.Context) {
	var req model.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request body")
		return
	}

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

	categoryId, err := category_service.CreateService(userStrId, req.ParentId, req.Name, req.Type)
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseSuccess(c, gin.H{
		"category_id": categoryId,
	})
}