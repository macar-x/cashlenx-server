package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

var (
	openapi      *openapi3.T
	routesRouter routers.Router
)

func init() {
	// Only initialize schema validation if enabled
	if util.GetConfigByKey("api.schema.validation") != "true" {
		return
	}

	// Load OpenAPI spec from file
	specPath := "docs/openapi.yaml"
	data, err := os.ReadFile(specPath)
	if err != nil {
		panic("Failed to load OpenAPI spec: " + err.Error())
	}

	// Parse OpenAPI spec
	spec, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		panic("Failed to parse OpenAPI spec: " + err.Error())
	}

	// Validate OpenAPI spec
	if err := spec.Validate(context.Background()); err != nil {
		panic("Invalid OpenAPI spec: " + err.Error())
	}

	// Create router for path matching
	openapi = spec
	routesRouter, err = gorillamux.NewRouter(openapi)
	if err != nil {
		panic("Failed to create OpenAPI router: " + err.Error())
	}
}

// SchemaValidation middleware to validate requests against OpenAPI schema
func SchemaValidation(next http.Handler) http.Handler {
	// Return original handler if schema validation is disabled
	if util.GetConfigByKey("api.schema.validation") != "true" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip health and version endpoints
		if strings.HasPrefix(r.URL.Path, "/api/health") || strings.HasPrefix(r.URL.Path, "/api/version") {
			next.ServeHTTP(w, r)
			return
		}

		// Validate request against schema
		if err := validateRequest(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.NewErrorResponse("VALIDATION_ERROR", err.Error())
			json.NewEncoder(w).Encode(response)
			return
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// validateRequest validates incoming request against OpenAPI schema
func validateRequest(r *http.Request) error {
	// Find matching route
	route, pathParams, err := routesRouter.FindRoute(r)
	if err != nil {
		return err
	}

	// Create validation context
	ctx := context.Background()

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	return openapi3filter.ValidateRequest(ctx, requestValidationInput)
}
