package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/controller/cash_flow_controller"
	"github.com/macar-x/cashlenx-server/controller/category_controller"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/middleware"
)

func StartServer(port int32) {
	r := mux.NewRouter()

	// Register routes
	registerHealthRoutes(r)
	registerCashRoute(r)
	registerCategoryRoute(r)

	// Apply middleware
	handler := middleware.Logging(middleware.CORS(r))

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("API server is running on http://localhost%s\n", addr)
	http.ListenAndServe(addr, handler)
}

func registerHealthRoutes(r *mux.Router) {
	r.HandleFunc("/api/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/version", versionInfo).Methods("GET")
}

func registerCashRoute(r *mux.Router) {
	// Create
	r.HandleFunc("/api/cash/outcome", cash_flow_controller.CreateOutcome).Methods("POST")
	r.HandleFunc("/api/cash/income", cash_flow_controller.CreateIncome).Methods("POST")

	// Read
	r.HandleFunc("/api/cash/list", cash_flow_controller.ListAll).Methods("GET")
	r.HandleFunc("/api/cash/{id}", cash_flow_controller.QueryById).Methods("GET")
	r.HandleFunc("/api/cash/date/{date}", cash_flow_controller.QueryByDate).Methods("GET")
	r.HandleFunc("/api/cash/range", cash_flow_controller.QueryByDateRange).Methods("GET")

	// Summary endpoints
	r.HandleFunc("/api/cash/summary/daily/{date}", cash_flow_controller.GetDailySummary).Methods("GET")
	r.HandleFunc("/api/cash/summary/monthly/{month}", cash_flow_controller.GetMonthlySummary).Methods("GET")
	r.HandleFunc("/api/cash/summary/yearly/{year}", cash_flow_controller.GetYearlySummary).Methods("GET")

	// Update
	r.HandleFunc("/api/cash/{id}", cash_flow_controller.UpdateById).Methods("PUT")

	// Delete
	r.HandleFunc("/api/cash/{id}", cash_flow_controller.DeleteById).Methods("DELETE")
	r.HandleFunc("/api/cash/date/{date}", cash_flow_controller.DeleteByDate).Methods("DELETE")
}

func registerCategoryRoute(r *mux.Router) {
	// Create
	r.HandleFunc("/api/category", category_controller.Create).Methods("POST")

	// Read
	r.HandleFunc("/api/category/list", category_controller.ListAll).Methods("GET")
	r.HandleFunc("/api/category/{id}", category_controller.QueryById).Methods("GET")
	r.HandleFunc("/api/category/name/{name}", category_controller.QueryByName).Methods("GET")
	r.HandleFunc("/api/category/children/{parent_id}", category_controller.QueryChildren).Methods("GET")

	// Update
	r.HandleFunc("/api/category/{id}", category_controller.UpdateById).Methods("PUT")

	// Delete
	r.HandleFunc("/api/category/{id}", category_controller.DeleteById).Methods("DELETE")
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "cashlenx-api",
		"message": "API is running",
	})
}

// Version info endpoint
func versionInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version":     model.Version,
		"name":        "CashLenX API",
		"description": "Personal finance management API",
		"endpoints": map[string][]string{
			"cash_flow": {
				"POST /api/cash/outcome",
				"POST /api/cash/income",
				"GET /api/cash/list",
				"GET /api/cash/{id}",
				"GET /api/cash/date/{date}",
				"GET /api/cash/range?from=YYYYMMDD&to=YYYYMMDD",
				"GET /api/cash/summary/daily/{date}",
				"GET /api/cash/summary/monthly/{month}",
				"GET /api/cash/summary/yearly/{year}",
				"PUT /api/cash/{id}",
				"DELETE /api/cash/{id}",
				"DELETE /api/cash/date/{date}",
			},
			"category": {
				"POST /api/category",
				"GET /api/category/list",
				"GET /api/category/{id}",
				"GET /api/category/name/{name}",
				"GET /api/category/children/{parent_id}",
				"PUT /api/category/{id}",
				"DELETE /api/category/{id}",
			},
			"health": {
				"GET /api/health",
				"GET /api/version",
			},
		},
	})
}
