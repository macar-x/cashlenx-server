package category_controller

import (
	"net/http"
	"strconv"

	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// ListAll returns paginated list of all categories
func ListAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Default limit for categories
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
	categories, totalCount, err := category_service.ListAllService(limit, offset)
	if err != nil {
		util.ComposeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return with pagination metadata
	util.ComposeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":        categories,
		"total_count": totalCount,
		"limit":       limit,
		"offset":      offset,
	})
}
