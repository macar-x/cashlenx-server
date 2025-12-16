package util

import (
	"encoding/json"
	"net/http"

	"github.com/macar-x/cashlenx-server/errors"
)

// ParseJSONRequest is a utility function to parse JSON requests
func ParseJSONRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	return err
}

// ErrorInfo defines the structure for error responses
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MetaInfo defines metadata structure for responses
type MetaInfo struct {
	Total        int64  `json:"total,omitempty"`
	Page         int64  `json:"page,omitempty"`
	Limit        int64  `json:"limit,omitempty"`
	Available    bool   `json:"available,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	ResponseTime int64  `json:"response_time,omitempty"`
}

// ResponseWrapper defines a consistent response structure for all API endpoints
type ResponseWrapper struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorInfo  `json:"error,omitempty"`
	Meta  *MetaInfo   `json:"meta,omitempty"`
}

// ComposeJSONResponse is a utility function to write JSON responses with consistent wrapper
func ComposeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response ResponseWrapper

	// Check if data is an error
	if err, ok := data.(error); ok {
		// Check if it's an AppError from errors package
		if appErr, ok := err.(*errors.AppError); ok {
			// Create error response from AppError
			response = ResponseWrapper{
				Error: &ErrorInfo{
					Code:    string(appErr.Code),
					Message: appErr.Message,
				},
			}
		} else {
			// Create generic error response
			response = ResponseWrapper{
				Error: &ErrorInfo{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			}
		}
	} else if errMap, ok := data.(map[string]string); ok {
		// Check if data is an error map
		if errMsg, hasError := errMap["error"]; hasError {
			// Create error response from error map
			response = ResponseWrapper{
				Error: &ErrorInfo{
					Code:    "BAD_REQUEST",
					Message: errMsg,
				},
			}
		} else {
			// Create success response with message map
			response = ResponseWrapper{
				Data: errMap,
			}
		}
	} else if msgMap, ok := data.(map[string]interface{}); ok {
		// Check if data is an error map interface
		if errMsg, hasError := msgMap["error"].(string); hasError {
			// Create error response from error map interface
			response = ResponseWrapper{
				Error: &ErrorInfo{
					Code:    "BAD_REQUEST",
					Message: errMsg,
				},
			}
		} else {
			// Create success response with data map
			response = ResponseWrapper{
				Data: msgMap,
			}
		}
	} else {
		// Create success response with raw data
		response = ResponseWrapper{
			Data: data,
		}
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}
