package category_controller

import (
	"encoding/json"
	"net/http"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "invalid request body",
		})
		return
	}

	// 从请求上下文中获取用户ID
	userId := r.Context().Value("user_id")
	if userId == nil {
		util.ComposeJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "user not authenticated",
		})
		return
	}

	// 验证用户ID格式
	userIdStr, ok := userId.(string)
	if !ok {
		util.ComposeJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "invalid user ID format",
		})
		return
	}

	categoryId, err := category_service.CreateService(userIdStr, req.ParentId, req.Name, req.Type)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"data": map[string]interface{}{
			"category_id": categoryId,
		},
	})
}
