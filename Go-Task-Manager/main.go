package main

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"go-task-manager/handlers"
	"go-task-manager/services"
	"go-task-manager/storage"
)

func main() {
	// Explicitly register standard MIME types to prevent Windows registry misconfigurations
	_ = mime.AddExtensionType(".css", "text/css; charset=utf-8")
	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".json", "application/json; charset=utf-8")
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	_ = mime.AddExtensionType(".png", "image/png")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize data storage
	dataPath := filepath.Join("data", "tasks.json")
	store := storage.NewJSONStorage(dataPath)

	// Initialize business logic service
	taskService := services.NewTaskService(store)

	// Load HTML templates
	tmplPattern := filepath.Join("templates", "*.html")
	tmpl, err := template.ParseGlob(tmplPattern)
	if err != nil {
		log.Printf("Warning: Failed to parse templates (%s): %v", tmplPattern, err)
	}

	// Initialize HTTP handlers
	taskHandler := handlers.NewTaskHandler(taskService, tmpl)

	// Create request multiplexer
	mux := http.NewServeMux()

	// Serve static assets (/static/css/, /static/js/, /static/assets/)
	staticDir := http.Dir("static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Web interface routes (Stitch imported multi-screen navigation)
	mux.HandleFunc("/", taskHandler.ServeHome)
	mux.HandleFunc("/dashboard", taskHandler.ServeHome)
	mux.HandleFunc("/tasks", taskHandler.ServeAllTasks)
	mux.HandleFunc("/empty", taskHandler.ServeEmptyState)

	// API routes
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTask(w, r)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	serverAddr := ":" + port
	fmt.Println("==================================================")
	fmt.Println(" TaskFlow – Personal Task Manager (Go Server)")
	fmt.Println("==================================================")
	fmt.Printf(" Server running at: http://localhost:%s/\n", port)
	fmt.Println(" Static assets served from: /static/")
	fmt.Println(" Initial setup and Go foundation ready.")
	fmt.Println("==================================================")

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
