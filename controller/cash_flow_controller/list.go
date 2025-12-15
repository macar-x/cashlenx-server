package cash_flow_controller

import (
	"net/http"
	"strconv"

	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/util"
)

// ListAll returns paginated list of all cash flows
func ListAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	cashType := r.URL.Query().Get("type") // Optional: INCOME or OUTCOME

	limit := 20 // Default limit
	offset := 0 // Default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Call service to get paginated results
	cashFlows, totalCount, err := cash_flow_service.QueryAll(cashType, limit, offset)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return with pagination metadata
	util.ComposeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":        cashFlows,
		"total_count": totalCount,
		"limit":       limit,
		"offset":      offset,
	})
}
