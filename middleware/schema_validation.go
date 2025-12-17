package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/macar-x/cashlenx-server/errors"
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
			// Use ComposeJSONResponse for consistent error formatting
			util.ComposeJSONResponse(w, http.StatusBadRequest, errors.NewValidationError("Request body validation failed. Please check your input format."))
			return
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// validateRequest validates incoming request against OpenAPI schema
func validateRequest(r *http.Request) error {
	// Create a copy of the request with modified URL to match OpenAPI server URL
	// This ensures validation works regardless of the actual hostname/port
	rCopy := r.Clone(context.Background())
	
	// Use the first server URL from the spec or default to http://localhost:8080
	var serverURL string
	if len(openapi.Servers) > 0 {
		serverURL = openapi.Servers[0].URL
	} else {
		serverURL = "http://localhost:8080"
	}
	
	// Parse the server URL
	server, err := r.URL.Parse(serverURL)
	if err != nil {
		return err
	}
	
	// Keep the original path, query, fragment, etc.
	rCopy.URL.Scheme = server.Scheme
	rCopy.URL.Host = server.Host
	// Don't change the path, query, fragment, etc.
	
	// Find matching route using the modified URL
	route, pathParams, err := routesRouter.FindRoute(rCopy)
	if err != nil {
		return err
	}

	// Create validation context
	ctx := context.Background()

	// Validate request using the original request but matched route
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	return openapi3filter.ValidateRequest(ctx, requestValidationInput)
}
