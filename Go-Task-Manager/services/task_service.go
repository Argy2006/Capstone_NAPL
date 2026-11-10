package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"time"

	"go-task-manager/models"
	"go-task-manager/storage"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTask  = errors.New("invalid task")
)

type TaskService struct {
	storage storage.Storage
	mu      sync.Mutex
}

func NewTaskService(s storage.Storage) *TaskService {
	return &TaskService{
		storage: s,
	}
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.storage.LoadTasks()
}

func (s *TaskService) GetTaskByID(id string) (*models.Task, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}

	return nil, ErrTaskNotFound
}

func (s *TaskService) CreateTask(task models.Task) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.Title = strings.TrimSpace(task.Title)

	if task.Title == "" {
		return nil, ErrInvalidTask
	}

	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}

	if task.Status == "" {
		task.Status = models.StatusPending
	}

	if !validPriority(task.Priority) || !validStatus(task.Status) {
		return nil, ErrInvalidTask
	}

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	task.ID = id
	task.CreatedAt = now
	task.UpdatedAt = now

	tasks = append(tasks, task)

	if err := s.storage.SaveTasks(tasks); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) UpdateTask(id string, updated models.Task) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updated.Title = strings.TrimSpace(updated.Title)

	if updated.Title == "" {
		return nil, ErrInvalidTask
	}

	if !validPriority(updated.Priority) || !validStatus(updated.Status) {
		return nil, ErrInvalidTask
	}

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			updated.ID = id
			updated.CreatedAt = tasks[i].CreatedAt
			updated.UpdatedAt = time.Now()

			tasks[i] = updated

			if err := s.storage.SaveTasks(tasks); err != nil {
				return nil, err
			}

			return &tasks[i], nil
		}
	}

	return nil, ErrTaskNotFound
}

func (s *TaskService) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.storage.SaveTasks(tasks)
		}
	}

	return ErrTaskNotFound
}

func validPriority(priority models.Priority) bool {
	return priority == models.PriorityLow ||
		priority == models.PriorityMedium ||
		priority == models.PriorityHigh
}

func validStatus(status models.Status) bool {
	return status == models.StatusPending ||
		status == models.StatusInProgress ||
		status == models.StatusCompleted
}

func generateID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
