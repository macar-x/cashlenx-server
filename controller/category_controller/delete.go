package category_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// DeleteById deletes a category by ID
func DeleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("id is required"))
		return
	}

	// Get the category before deleting to return it
	categoryToDelete, err := category_service.QueryService(id, "", "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(categoryToDelete) == 0 {
		util.ComposeJSONResponse(w, http.StatusNotFound, errors.NewAppError(errors.ErrNotFound, "category not found", nil))
		return
	}

	// Delete the category
	if err := category_service.DeleteService(id, ""); err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return the deleted category
	util.ComposeJSONResponse(w, http.StatusOK, categoryToDelete[0])
}
