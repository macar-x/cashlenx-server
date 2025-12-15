package model

// ResponseWrapper defines a consistent response structure for all API endpoints
type ResponseWrapper struct {
	Data  interface{} `json:"data,omitempty"`  // The actual data payload
	Error *ErrorInfo  `json:"error,omitempty"` // Error information if any
	Meta  *MetaInfo   `json:"meta,omitempty"`  // Metadata about the response
}

// ErrorInfo defines the structure for error responses
type ErrorInfo struct {
	Code    string `json:"code"`    // Error code for machine consumption
	Message string `json:"message"` // Human-readable error message
	Details string `json:"details,omitempty"` // Additional error details
}

// MetaInfo defines metadata structure for responses
type MetaInfo struct {
	Total       int64 `json:"total,omitempty"`       // Total number of items (for pagination)
	Page        int64 `json:"page,omitempty"`        // Current page number
	Limit       int64 `json:"limit,omitempty"`       // Items per page
	Available   bool  `json:"available,omitempty"`   // Resource availability flag
	RequestID   string `json:"request_id,omitempty"`  // Unique request identifier
	ResponseTime int64 `json:"response_time,omitempty"` // Response time in milliseconds
}

// NewSuccessResponse creates a successful response with data
func NewSuccessResponse(data interface{}) ResponseWrapper {
	return ResponseWrapper{
		Data: data,
	}
}

// NewSuccessResponseWithMeta creates a successful response with data and metadata
func NewSuccessResponseWithMeta(data interface{}, meta *MetaInfo) ResponseWrapper {
	return ResponseWrapper{
		Data: data,
		Meta: meta,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(code, message string) ResponseWrapper {
	return ResponseWrapper{
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// NewErrorResponseWithDetails creates an error response with additional details
func NewErrorResponseWithDetails(code, message, details string) ResponseWrapper {
	return ResponseWrapper{
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
