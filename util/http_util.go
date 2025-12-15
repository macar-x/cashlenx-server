package util

import (
	"encoding/json"
	"net/http"
)

// ResponseWrapper defines a consistent response structure for all API endpoints
type ResponseWrapper struct {
	// Data contains the actual response payload
	Data  interface{} `json:"data,omitempty"`
	// Error contains error information if the request failed
	Error *ErrorInfo  `json:"error,omitempty"`
	// Meta contains metadata about the response
	Meta  *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo defines the structure for error responses
type ErrorInfo struct {
	// Code is the error code for machine consumption
	Code    string `json:"code"`
	// Message is the human-readable error message
	Message string `json:"message"`
}

// MetaInfo defines metadata structure for responses
type MetaInfo struct {
	// Total is the total number of items (for pagination)
	Total int64 `json:"total,omitempty"`
	// Page is the current page number
	Page  int64 `json:"page,omitempty"`
	// Limit is the number of items per page
	Limit int64 `json:"limit,omitempty"`
}

// ParseJSONRequest is a utility function to parse JSON requests
func ParseJSONRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	return err
}

// ComposeJSONResponse is a utility function to write JSON responses with consistent wrapper
func ComposeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response ResponseWrapper

	// Check if data is an error
	if err, ok := data.(error); ok {
		// Check if it's an AppError from errors package
		if appErr, ok := err.(interface{
			GetCode() string
			GetMessage() string
		}); ok {
			// Create error response from AppError interface
			response = ResponseWrapper{
				Error: &ErrorInfo{
					Code:    appErr.GetCode(),
					Message: appErr.GetMessage(),
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
