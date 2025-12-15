package category_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// DeleteById deletes a category by ID
func DeleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plainId := vars["id"]

	if plainId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	err := category_service.DeleteService(plainId, "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"message": "category deleted successfully"})
}
