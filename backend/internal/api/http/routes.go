package http

import (
	"log"
	"net/http"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/http/handlers"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes() error {
	// Create diagram handler
	diagramHandler, err := handlers.NewDiagramHandler()
	if err != nil {
		return err
	}

	// API routes
	http.HandleFunc("/api/diagrams/process", diagramHandler.ProcessDiagram)
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/", RootHandler)

	return nil
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"arch-visualizer-api"}`))
}

// RootHandler handles root requests
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Arch Visualizer API","version":"1.0.0"}`))
}

// StartServer starts the HTTP server
func StartServer(port string) error {
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("API endpoints:")
	log.Printf("  POST /api/diagrams/process - Process a diagram")
	log.Printf("  GET  /api/health - Health check")
	log.Printf("  GET  / - Root endpoint")

	return http.ListenAndServe(":"+port, nil)
}
