package cash_flow_controller

import (
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

// QueryByDateRange queries cash flows between two dates
func QueryByDateRange(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	if fromDate == "" || toDate == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewInvalidInputError("from and to dates are required"))
		return
	}

	// Call service to get records in range
	cashFlowEntities, err := cash_flow_service.QueryByDateRange(fromDate, toDate)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, cashFlowEntities)
}
