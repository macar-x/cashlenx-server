package cash_flow_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

func DeleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plainId := vars["id"]
	if plainId == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is empty"})
	}
	cashFlowEntity, err := cash_flow_service.DeleteById(plainId)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntity)
}

func DeleteByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	if date == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "date is empty"})
	}
	cashFlowEntityList, err := cash_flow_service.DeleteByDate(date)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntityList)
}
