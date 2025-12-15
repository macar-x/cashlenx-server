package cash_flow_controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

// UpdateById updates a cash flow record by ID
func UpdateById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plainId := vars["id"]
	if plainId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	// Parse JSON body for update fields
	var requestBody map[string]interface{}
	if err := util.ParseJSONRequest(r, &requestBody); err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Extract optional fields
	belongsDate, _ := requestBody["belongs_date"].(string)
	categoryName, _ := requestBody["category_name"].(string)
	description, _ := requestBody["description"].(string)

	var amount float64
	if amountVal, ok := requestBody["amount"]; ok {
		switch v := amountVal.(type) {
		case float64:
			amount = v
		case string:
			amount, _ = strconv.ParseFloat(v, 64)
		}
	}

	// Call service to update
	updatedEntity, err := cash_flow_service.UpdateById(plainId, belongsDate, categoryName, amount, description)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, updatedEntity)
}
