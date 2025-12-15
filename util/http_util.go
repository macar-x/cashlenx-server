package util

import (
	"encoding/json"
	"net/http"
)

// ParseJSONRequest is a utility function to parse JSON requests
func ParseJSONRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	return err
}

// JSONResponse is a utility function to write JSON responses
func ComposeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode the data as JSON and write it to the response writer
	json.NewEncoder(w).Encode(data)
}
