package category_controller

import (
	"net/http"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// Create creates a new category
func Create(w http.ResponseWriter, r *http.Request) {
	var requestBody model.CategoryDTO
	if err := util.ParseJSONRequest(r, &requestBody); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if requestBody.Name == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "category name is required"})
		return
	}

	plainId, err := category_service.CreateService(requestBody.ParentName, requestBody.Name)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, map[string]string{
		"id":      plainId,
		"message": "category created successfully",
	})
}
