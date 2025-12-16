package category_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// UpdateById updates a category by ID
func UpdateById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plainId := vars["id"]

	if plainId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("id is required"))
		return
	}

	// Parse JSON body for update fields
	var requestBody map[string]interface{}
	if err := util.ParseJSONRequest(r, &requestBody); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	// Extract optional fields
	parentPlainId, _ := requestBody["parent_id"].(string)
	categoryName, _ := requestBody["name"].(string)
	categoryType, _ := requestBody["type"].(string)

	// Call service to update
	err := category_service.UpdateService(plainId, parentPlainId, categoryName, categoryType)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Get the updated category entity
	updatedCategory, err := category_service.QueryService(plainId, "", "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(updatedCategory) == 0 {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError("failed to retrieve updated category", nil))
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, updatedCategory[0])
}
