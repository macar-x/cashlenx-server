package statistic_controller

import (
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
	"github.com/macar-x/cashlenx-server/service/statistic_service"
	"github.com/macar-x/cashlenx-server/util"
)

// ExportData exports user's cash flow data to Excel file
func ExportData(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request context
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		util.ComposeJSONResponse(w, http.StatusUnauthorized, errors.NewUnauthorizedError("user not authenticated"))
		return
	}

	// Parse query parameters
	fromDate := r.URL.Query().Get("from_date") // YYYYMMDD or YYYY-MM-DD
	toDate := r.URL.Query().Get("to_date")     // YYYYMMDD or YYYY-MM-DD
	filePath := r.URL.Query().Get("file_path")

	// Default file path if not provided
	if filePath == "" {
		filePath = "./export.xlsx"
	}

	// Call service to export data for this user
	err := statistic_service.ExportForUser(fromDate, toDate, filePath, userId)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError(err.Error()))
		return
	}

	response := map[string]interface{}{
		"message":   "Data exported successfully",
		"file_path": filePath,
		"user_id":   userId,
	}
	util.ComposeJSONResponse(w, http.StatusOK, response)
}
