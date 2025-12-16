package model

// NewSuccessResponse creates a successful response with data
func NewSuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
	}
}

// NewSuccessResponseWithMeta creates a successful response with data and metadata
func NewSuccessResponseWithMeta(data interface{}, meta map[string]interface{}) map[string]interface{} {
	response := map[string]interface{}{
		"data": data,
	}
	if meta != nil {
		response["meta"] = meta
	}
	return response
}

// NewErrorResponse creates an error response
func NewErrorResponse(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}
}

// NewErrorResponseWithDetails creates an error response with additional details
func NewErrorResponseWithDetails(code, message, details string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
			"details": details,
		},
	}
}
