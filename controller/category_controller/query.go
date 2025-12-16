package category_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// QueryById retrieves a category by ID
func QueryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("id is required"))
		return
	}

	categories, err := category_service.QueryService(id, "", "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(categories) == 0 {
		util.ComposeJSONResponse(w, http.StatusNotFound, errors.NewAppError(errors.ErrNotFound, "category not found", nil))
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categories[0])
}

// QueryByName retrieves categories by name (fuzzy match)
func QueryByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("name is required"))
		return
	}

	categories, err := category_service.QueryService("", "", name)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categories)
}

// QueryChildren retrieves children categories by parent ID
func QueryChildren(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentId := vars["parent_id"]

	if parentId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("parent_id is required"))
		return
	}

	categories, err := category_service.QueryService("", parentId, "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categories)
}
