package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"go-task-manager/models"
	"go-task-manager/services"
)

type TaskHandler struct {
	service   *services.TaskService
	templates *template.Template
}

func NewTaskHandler(service *services.TaskService, tmpl *template.Template) *TaskHandler {
	return &TaskHandler{
		service:   service,
		templates: tmpl,
	}
}

func (h *TaskHandler) renderView(w http.ResponseWriter, tmplName string, data map[string]interface{}) {
	if h.templates == nil {
		http.Error(w, "Templates are not available", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := h.templates.ExecuteTemplate(buf, tmplName, data); err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (h *TaskHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title": "TaskFlow — Dashboard Overview & Task Management",
	}

	h.renderView(w, "index.html", data)
}

func (h *TaskHandler) ServeAllTasks(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "TaskFlow — All Tasks Workspace & Table View",
	}

	h.renderView(w, "all_tasks.html", data)
}

func (h *TaskHandler) ServeEmptyState(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "TaskFlow — Empty State & Getting Started",
	}

	h.renderView(w, "empty_state.html", data)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func decodeTask(r *http.Request) (models.Task, error) {
	var task models.Task

	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		return task, err
	}

	return task, nil
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch tasks",
		})
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID",
		})
		return
	}

	task, err := h.service.GetTaskByID(id)
	if errors.Is(err, services.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch task",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid task data",
		})
		return
	}

	created, err := h.service.CreateTask(task)
	if errors.Is(err, services.ErrInvalidTask) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Task title, priority or status is invalid",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create task",
		})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID",
		})
		return
	}

	task, err := decodeTask(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid task data",
		})
		return
	}

	updated, err := h.service.UpdateTask(id, task)
	if errors.Is(err, services.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}
	if errors.Is(err, services.ErrInvalidTask) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Task title, priority or status is invalid",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update task",
		})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID",
		})
		return
	}

	err := h.service.DeleteTask(id)
	if errors.Is(err, services.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete task",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Task deleted successfully",
	})
}
