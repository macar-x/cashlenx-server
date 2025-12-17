package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/macar-x/cashlenx-server/controller/auth_controller"
	"github.com/macar-x/cashlenx-server/controller/cash_flow_controller"
	"github.com/macar-x/cashlenx-server/controller/category_controller"
	"github.com/macar-x/cashlenx-server/controller/manage_controller"
	"github.com/macar-x/cashlenx-server/controller/user_controller"
	"github.com/macar-x/cashlenx-server/middleware"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

func StartServer(port int32) {
	// Explicitly load timezone at server startup to ensure it's configured
	// and logged immediately
	tz := util.GetTimezone()
	fmt.Printf("Loaded timezone: %v\n", tz)

	r := mux.NewRouter()

	// Register routes
	registerHealthRoutes(r)
	registerUserRoute(r)
	registerCashRoute(r)
	registerCategoryRoute(r)
	registerManageRoute(r)

	// Apply middleware
	handler := middleware.Logging(middleware.Auth(middleware.SchemaValidation(middleware.CORS(r))))

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("API server is running on http://localhost%s\n", addr)
	http.ListenAndServe(addr, handler)
}

func registerHealthRoutes(r *mux.Router) {
	r.HandleFunc("/api/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/version", versionInfo).Methods("GET")
}

func registerUserRoute(r *mux.Router) {
	// User management routes
	r.HandleFunc("/api/user", user_controller.Create).Methods("POST")
	r.HandleFunc("/api/user", user_controller.ListAll).Methods("GET")
	r.HandleFunc("/api/user/{id}", user_controller.Get).Methods("GET")
	r.HandleFunc("/api/user/{id}", user_controller.Update).Methods("PUT")
	r.HandleFunc("/api/user/{id}", user_controller.Delete).Methods("DELETE")

	// Authentication routes
	r.HandleFunc("/api/auth/login", auth_controller.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", auth_controller.Register).Methods("POST")
}

func registerCashRoute(r *mux.Router) {
	// Create
	r.HandleFunc("/api/cash/expense", cash_flow_controller.CreateExpense).Methods("POST")
	r.HandleFunc("/api/cash/income", cash_flow_controller.CreateIncome).Methods("POST")

	// Read
	r.HandleFunc("/api/cash", cash_flow_controller.ListAll).Methods("GET")
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
	// Read all categories with filtering
	r.HandleFunc("/api/category", category_controller.ListAll).Methods("GET")
	// Read specific category
	r.HandleFunc("/api/category/{id}", category_controller.QueryById).Methods("GET")
	// Read by name
	r.HandleFunc("/api/category/name/{name}", category_controller.QueryByName).Methods("GET")
	// Read children categories - RESTful design: parent/{id}/children
	r.HandleFunc("/api/category/{parent_id}/children", category_controller.QueryChildren).Methods("GET")
	// Read category tree structure
	r.HandleFunc("/api/category/tree", category_controller.Tree).Methods("GET")

	// Update
	r.HandleFunc("/api/category/{id}", category_controller.UpdateById).Methods("PUT")

	// Delete
	r.HandleFunc("/api/category/{id}", category_controller.DeleteById).Methods("DELETE")
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "cashlenx-api",
		"message": "API is running",
	}
	util.ComposeJSONResponse(w, http.StatusOK, response)
}

func registerManageRoute(r *mux.Router) {
	// Dump and restore endpoints
	r.HandleFunc("/api/manage/dump", manage_controller.DumpDatabase).Methods("GET")
	r.HandleFunc("/api/manage/restore", manage_controller.RestoreDatabase).Methods("POST")

	// Import and export endpoints
	r.HandleFunc("/api/manage/export", manage_controller.ExportData).Methods("GET")
	r.HandleFunc("/api/manage/import", manage_controller.ImportData).Methods("POST")
}

// Version info endpoint
func versionInfo(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"version":     model.Version,
		"name":        "CashLenX API",
		"description": "Personal finance management API",
		"endpoints": map[string][]string{
			"cash_flow": {
				"POST /api/cash/expense",
				"POST /api/cash/income",
				"GET /api/cash",
				"GET /api/cash?limit=20&offset=0&type=income",
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
				"GET /api/category",
				"GET /api/category?type=income&parent_id=XXX",
				"GET /api/category/{id}",
				"GET /api/category/name/{name}",
				"GET /api/category/{parent_id}/children",
				"PUT /api/category/{id}",
				"DELETE /api/category/{id}",
			},
			"manage": {
				"GET /api/manage/dump",
				"POST /api/manage/restore",
				"GET /api/manage/export",
				"POST /api/manage/import",
			},
			"user": {
				"POST /api/user",
				"GET /api/user",
				"GET /api/user/{id}",
				"PUT /api/user/{id}",
				"DELETE /api/user/{id}",
			},
			"auth": {
				"POST /api/auth/login",
				"POST /api/auth/register",
			},
			"health": {
				"GET /api/health",
				"GET /api/version",
			},
		},
	}
	util.ComposeJSONResponse(w, http.StatusOK, response)
}
