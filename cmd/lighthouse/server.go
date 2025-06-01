package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// serverCmd runs a lighthouse server
func serverCmd() *cobra.Command {
	var (
		port    int
		dataDir string
		tlsCert string
		tlsKey  string
	)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run a Lighthouse server",
		Long:  "Start a server to coordinate pledges for projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(port, dataDir, tlsCert, tlsKey)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	cmd.Flags().StringVarP(&dataDir, "data", "d", "./lighthouse-data", "Data directory for projects and pledges")
	cmd.Flags().StringVar(&tlsCert, "tls-cert", "", "TLS certificate file")
	cmd.Flags().StringVar(&tlsKey, "tls-key", "", "TLS key file")

	return cmd
}

func runServer(port int, dataDir, tlsCert, tlsKey string) error {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	fmt.Printf("Starting Lighthouse server on port %d\n", port)
	fmt.Printf("Data directory: %s\n", dataDir)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Project routes
	mux.HandleFunc("/api/projects", corsMiddleware(projectsHandler(dataDir)))
	mux.HandleFunc("/api/projects/", corsMiddleware(projectHandler(dataDir)))

	// Pledge routes
	mux.HandleFunc("/api/pledges", corsMiddleware(pledgesHandler(dataDir)))

	// Add logging middleware
	handler := loggingMiddleware(mux)

	// Start server
	addr := fmt.Sprintf(":%d", port)

	if tlsCert != "" && tlsKey != "" {
		fmt.Printf("Starting HTTPS server on %s\n", addr)
		return http.ListenAndServeTLS(addr, tlsCert, tlsKey, handler)
	} else {
		fmt.Printf("Starting HTTP server on %s\n", addr)
		return http.ListenAndServe(addr, handler)
	}
}

// Middleware for CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "lighthouse-server",
	})
}

// Projects handler
func projectsHandler(dataDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			// List all projects
			projects, err := listProjects(dataDir)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to list projects: %v", err), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"projects": projects})

		case "POST":
			// Create new project (placeholder)
			w.WriteHeader(http.StatusNotImplemented)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Project creation via API not yet implemented",
				"note":  "Use CLI: lighthouse project create",
			})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Individual project handler
func projectHandler(dataDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract project ID from URL
		projectID := filepath.Base(r.URL.Path)

		switch r.Method {
		case "GET":
			// Get project details (placeholder)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"project": map[string]string{
					"id":      projectID,
					"message": "Project details not yet implemented",
				},
			})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Pledges handler
func pledgesHandler(dataDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			// List pledges (placeholder)
			json.NewEncoder(w).Encode(map[string]interface{}{"pledges": []string{}})

		case "POST":
			// Create pledge (placeholder)
			w.WriteHeader(http.StatusNotImplemented)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Pledge creation via API not yet implemented",
				"note":  "Use CLI: lighthouse pledge create",
			})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// List projects in data directory
func listProjects(dataDir string) ([]map[string]interface{}, error) {
	pattern := filepath.Join(dataDir, "*.lighthouse")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var projects []map[string]interface{}
	for _, file := range files {
		// For now, just return basic file info
		// In a full implementation, we'd parse the project files
		base := filepath.Base(file)
		projects = append(projects, map[string]interface{}{
			"id":   base,
			"file": file,
			"name": base[:len(base)-11], // Remove .lighthouse extension
		})
	}

	return projects, nil
}
