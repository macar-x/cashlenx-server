package category_controller

import (
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// Create creates a new category
func Create(w http.ResponseWriter, r *http.Request) {
	var requestBody model.CategoryDTO
	if err := util.ParseJSONRequest(r, &requestBody); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	if requestBody.Name == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewValidationError("category name is required"))
		return
	}

	if requestBody.Type == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewValidationError("category type is required"))
		return
	}

	plainId, err := category_service.CreateService(requestBody.ParentId, requestBody.Name, requestBody.Type)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Get the created category entity
	createdCategory, err := category_service.QueryService(plainId, "", "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(createdCategory) == 0 {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError("failed to retrieve created category", nil))
		return
	}

	util.ComposeJSONResponse(w, http.StatusCreated, createdCategory[0])
}
