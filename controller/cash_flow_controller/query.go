package cash_flow_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

func QueryById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("id is required"))
		return
	}

	cashFlowEntity, err := cash_flow_service.QueryById(id)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntity)
}

func QueryByDate(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	if date == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("date is required"))
		return
	}

	cashFlowEntityList, err := cash_flow_service.QueryByDate(date)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntityList)
}
