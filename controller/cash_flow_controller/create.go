package cash_flow_controller

import (
	"errors"
	"net/http"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

func CreateOutcome(w http.ResponseWriter, r *http.Request) {
	requestBody, err := validCashFlowDTO(r)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	cashFlowEntity, err := cash_flow_service.SaveOutcome(requestBody.BelongsDate, requestBody.CategoryName, requestBody.Amount, requestBody.Description)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntity)
}

func CreateIncome(w http.ResponseWriter, r *http.Request) {
	requestBody, err := validCashFlowDTO(r)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	cashFlowEntity, err := cash_flow_service.SaveIncome(requestBody.BelongsDate, requestBody.CategoryName, requestBody.Amount, requestBody.Description)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntity)
}

func validCashFlowDTO(r *http.Request) (model.CashFlowDTO, error) {
	var requestBody model.CashFlowDTO
	err := util.ParseJSONRequest(r, &requestBody)
	if err != nil {
		return model.CashFlowDTO{}, err
	}

	if !cash_flow_service.IsOutcomeRequiredFiledSatisfied(requestBody.CategoryName, requestBody.Amount) {
		return model.CashFlowDTO{}, errors.New("some required fields are empty")
	}
	return requestBody, nil
}
