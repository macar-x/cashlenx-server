package category_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// QueryById queries a category by ID
func QueryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plainId := vars["id"]

	if plainId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	categoryEntities, err := category_service.QueryService(plainId, "", "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if len(categoryEntities) == 0 {
		util.ComposeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "category not found"})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categoryEntities[0])
}

// QueryByName queries categories by name
func QueryByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	categoryEntities, err := category_service.QueryService("", name, "")
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categoryEntities)
}

// QueryChildren queries child categories by parent ID
func QueryChildren(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentId := vars["parent_id"]

	if parentId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "parent_id is required"})
		return
	}

	categoryEntities, err := category_service.QueryService("", "", parentId)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, categoryEntities)
}
