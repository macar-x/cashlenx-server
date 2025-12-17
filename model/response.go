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
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

// NewErrorResponseWithDetails creates an error response with additional details
func NewErrorResponseWithDetails(code, message, details string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": details,
		},
	}
}

// NewFieldErrorResponse creates an error response for field validation
func NewFieldErrorResponse(code, message, field string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"field":   field,
		},
	}
}

// NewFieldErrorResponseWithDetails creates an error response for field validation with details
func NewFieldErrorResponseWithDetails(code, message, field, details string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"field":   field,
			"details": details,
		},
	}
}
