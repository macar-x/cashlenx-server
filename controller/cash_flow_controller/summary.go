package cash_flow_controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

// GetDailySummary returns summary for a specific day
func GetDailySummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	if date == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "date is required"})
		return
	}

	summary, err := cash_flow_service.GetSummary("daily", date)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, summary)
}

// GetMonthlySummary returns summary for a specific month (YYYYMM format)
func GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	month := vars["month"]

	if month == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "month is required (YYYYMM format)"})
		return
	}

	summary, err := cash_flow_service.GetSummary("monthly", month)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, summary)
}

// GetYearlySummary returns summary for a specific year (YYYY format)
func GetYearlySummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]

	if year == "" {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "year is required (YYYY format)"})
		return
	}

	summary, err := cash_flow_service.GetSummary("yearly", year)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	util.ComposeJSONResponse(w, http.StatusOK, summary)
}
