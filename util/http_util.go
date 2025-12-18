package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
)

// ParseJSONRequest is a utility function to parse JSON requests
func ParseJSONRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	return err
}

// Response defines the unified structure for all API responses
type Response struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Errors  map[string]string      `json:"errors,omitempty"`
}

// ComposeJSONResponse is a utility function to write JSON responses with consistent format
func ComposeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response Response

	// Initialize default values
	response.Code = "OK"
	response.Message = ""
	response.Data = nil

	// Check if data is an error
	if err, ok := data.(error); ok {
		// Check if it's an AppError from errors package
		if appErr, ok := err.(*errors.AppError); ok {
			// Set error code and message
			response.Code = string(appErr.Code)
			response.Message = "Validation failed"
			
			// Initialize errors map
			response.Errors = make(map[string]string)
			
			// Add field error
			if appErr.Field != "" {
				response.Errors[appErr.Field] = appErr.Message
			} else {
				// Generic error if no field specified
				response.Message = appErr.Message
			}
		} else {
			// Create generic error response for standard errors
			response.Code = "INTERNAL_ERROR"
			response.Message = err.Error()
		}
	} else {
		// Success response
		response.Data = data
		
		// Set appropriate message based on data type
		if user, ok := data.(map[string]interface{}); ok {
			if username, ok := user["username"].(string); ok {
				response.Message = "user " + username + " created"
			}
		} else if users, ok := data.([]interface{}); ok {
			if len(users) > 0 {
				response.Message = "users retrieved successfully"
			} else {
				response.Message = "no users found"
			}
		} else {
			response.Message = "operation successful"
		}
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

// SendFile sends a file as an HTTP response
func SendFile(w http.ResponseWriter, file io.Reader) {
	if _, err := io.Copy(w, file); err != nil {
		// If there's an error while sending the file, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError("Failed to send file", err))
	}
}
