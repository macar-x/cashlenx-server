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

// ErrorInfo defines the structure for error responses
type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Field     string `json:"field,omitempty"`
	RequestID string `json:"request_id,omitempty"`
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
			// Create detailed error response from AppError
			errorInfo := &ErrorInfo{
				Code:    string(appErr.Code),
				Message: appErr.Message,
				Field:   appErr.Field,
			}
			
			// Add cause details if available
			if appErr.Cause != nil {
				errorInfo.Details = appErr.Cause.Error()
			}
			
			response = ResponseWrapper{
				Error: errorInfo,
			}
		} else {
			// Create generic error response for standard errors
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
		if errObj, hasError := msgMap["error"].(map[string]interface{}); hasError {
			// Create detailed error response from error map
			errorInfo := &ErrorInfo{
				Code:    errObj["code"].(string),
				Message: errObj["message"].(string),
			}
			// Add optional fields if present
			if details, ok := errObj["details"].(string); ok {
				errorInfo.Details = details
			}
			if field, ok := errObj["field"].(string); ok {
				errorInfo.Field = field
			}
			response = ResponseWrapper{Error: errorInfo}
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

// SendFile sends a file as an HTTP response
func SendFile(w http.ResponseWriter, file io.Reader) {
	if _, err := io.Copy(w, file); err != nil {
		// If there's an error while sending the file, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		ComposeJSONResponse(w, http.StatusInternalServerError, errors.NewInternalError("Failed to send file", err))
	}
}
