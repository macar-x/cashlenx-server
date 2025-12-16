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

	if err := category_service.DeleteService(id, ""); err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"message": "category deleted successfully"})
}
