package tests

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-task-manager/handlers"
	"go-task-manager/models"
	"go-task-manager/services"
	"go-task-manager/storage"
)

func TestTaskModel(t *testing.T) {
	now := time.Now()
	task := models.Task{
		ID:          "task-1",
		Title:       "Test Task Creation",
		Description: "Initial test to verify Task model struct fields.",
		Priority:    models.PriorityHigh,
		Status:      models.StatusPending,
		Category:    "Academic",
		DueDate:     "2026-11-15",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if task.ID != "task-1" {
		t.Errorf("Expected ID 'task-1', got '%s'", task.ID)
	}

	if task.Priority != models.PriorityHigh {
		t.Errorf("Expected priority '%s', got '%s'", models.PriorityHigh, task.Priority)
	}

	if task.Status != models.StatusPending {
		t.Errorf("Expected status '%s', got '%s'", models.StatusPending, task.Status)
	}
}

func TestJSONStorage_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test_tasks.json")

	store := storage.NewJSONStorage(testFilePath)

	// Verify loading from non-existent file returns empty slice without error
	tasks, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed on new file: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(tasks))
	}

	// Save sample tasks
	sampleTasks := []models.Task{
		{
			ID:          "1",
			Title:       "Capstone Project Setup",
			Description: "Setup project structure and UI integration",
			Priority:    models.PriorityHigh,
			Status:      models.StatusInProgress,
			Category:    "SEC042",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := store.SaveTasks(sampleTasks); err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	// Verify persistence
	loaded, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed after save: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(loaded))
	}
	if loaded[0].Title != "Capstone Project Setup" {
		t.Errorf("Expected title 'Capstone Project Setup', got '%s'", loaded[0].Title)
	}
}

func TestTaskServiceInitialization(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "service_tasks.json")
	store := storage.NewJSONStorage(testFilePath)
	service := services.NewTaskService(store)

	tasks, err := service.GetAllTasks()
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(tasks))
	}
}

func TestHandler_ServeViews(t *testing.T) {
	// Parse all actual templates from parent directory
	tmplPattern := filepath.Join("..", "templates", "*.html")
	tmpl, err := template.ParseGlob(tmplPattern)
	if err != nil {
		t.Fatalf("Failed to parse templates at %s: %v", tmplPattern, err)
	}

	tmpDir := t.TempDir()
	store := storage.NewJSONStorage(filepath.Join(tmpDir, "tasks.json"))
	service := services.NewTaskService(store)
	handler := handlers.NewTaskHandler(service, tmpl)

	// Test 1: ServeHome (Dashboard)
	reqHome := httptest.NewRequest(http.MethodGet, "/", nil)
	wHome := httptest.NewRecorder()
	handler.ServeHome(wHome, reqHome)
	respHome := wHome.Result()
	defer respHome.Body.Close()

	if respHome.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK for Home, got %d", respHome.StatusCode)
	}
	if ct := respHome.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8' for Home, got '%s'", ct)
	}
	bodyHome := wHome.Body.String()
	if strings.HasPrefix(bodyHome, "\xef\xbb\xbf") {
		t.Errorf("Response body contains unexpected UTF-8 BOM prefix")
	}
	if !strings.HasPrefix(strings.TrimSpace(bodyHome), "<!DOCTYPE html>") {
		t.Errorf("Response body does not start with '<!DOCTYPE html>'")
	}
	if !strings.Contains(bodyHome, "TaskFlow") {
		t.Errorf("Expected home body to contain 'TaskFlow'")
	}
	if !strings.Contains(bodyHome, "/static/css/style.css") {
		t.Errorf("Expected home body to link to '/static/css/style.css'")
	}
	if !strings.Contains(bodyHome, "/static/js/app.js") {
		t.Errorf("Expected home body to link to '/static/js/app.js'")
	}

	// Test 2: ServeAllTasks (Table Workspace)
	reqTasks := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	wTasks := httptest.NewRecorder()
	handler.ServeAllTasks(wTasks, reqTasks)
	respTasks := wTasks.Result()
	defer respTasks.Body.Close()

	if respTasks.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK for All Tasks, got %d", respTasks.StatusCode)
	}
	if ct := respTasks.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8' for All Tasks, got '%s'", ct)
	}
	bodyTasks := wTasks.Body.String()
	if !strings.Contains(bodyTasks, "TaskFlow") {
		t.Errorf("Expected all tasks body to contain 'TaskFlow'")
	}

	// Test 3: ServeEmptyState
	reqEmpty := httptest.NewRequest(http.MethodGet, "/empty", nil)
	wEmpty := httptest.NewRecorder()
	handler.ServeEmptyState(wEmpty, reqEmpty)
	respEmpty := wEmpty.Result()
	defer respEmpty.Body.Close()

	if respEmpty.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK for Empty State, got %d", respEmpty.StatusCode)
	}
	if ct := respEmpty.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8' for Empty State, got '%s'", ct)
	}
	bodyEmpty := wEmpty.Body.String()
	if !strings.Contains(bodyEmpty, "TaskFlow") {
		t.Errorf("Expected empty state body to contain 'TaskFlow'")
	}
}

func TestHandler_GetTasks(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewJSONStorage(filepath.Join(tmpDir, "tasks.json"))
	service := services.NewTaskService(store)
	handler := handlers.NewTaskHandler(service, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetTasks(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json header, got %s", contentType)
	}

	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("Expected empty JSON array '[]', got '%s'", body)
	}
}

func TestStaticFilesExist(t *testing.T) {
	cssPath := filepath.Join("..", "static", "css", "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist", cssPath)
	}

	jsPath := filepath.Join("..", "static", "js", "app.js")
	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist", jsPath)
	}

	dataPath := filepath.Join("..", "data", "tasks.json")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist", dataPath)
	}
}

func TestStaticMimeTypes(t *testing.T) {
	mux := http.NewServeMux()
	staticDir := http.Dir(filepath.Join("..", "static"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	reqCss := httptest.NewRequest(http.MethodGet, "/static/css/style.css", nil)
	wCss := httptest.NewRecorder()
	mux.ServeHTTP(wCss, reqCss)
	respCss := wCss.Result()
	defer respCss.Body.Close()

	if respCss.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for /static/css/style.css, got %d", respCss.StatusCode)
	}
	if ct := respCss.Header.Get("Content-Type"); !strings.Contains(ct, "text/css") {
		t.Errorf("Expected Content-Type text/css, got %q", ct)
	}

	reqJs := httptest.NewRequest(http.MethodGet, "/static/js/app.js", nil)
	wJs := httptest.NewRecorder()
	mux.ServeHTTP(wJs, reqJs)
	respJs := wJs.Result()
	defer respJs.Body.Close()

	if respJs.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for /static/js/app.js, got %d", respJs.StatusCode)
	}
	if ct := respJs.Header.Get("Content-Type"); !strings.Contains(ct, "javascript") {
		t.Errorf("Expected Content-Type javascript, got %q", ct)
	}
}
