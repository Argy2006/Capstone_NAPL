package storage

import (
	"encoding/json"
	"os"
	"sync"

	"go-task-manager/models"
)

// Storage defines the persistence interface for task data.
type Storage interface {
	LoadTasks() ([]models.Task, error)
	SaveTasks(tasks []models.Task) error
}

// JSONStorage implements Storage using a local JSON file.
type JSONStorage struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONStorage creates a new instance of JSONStorage.
func NewJSONStorage(filePath string) *JSONStorage {
	return &JSONStorage{
		filePath: filePath,
	}
}

// LoadTasks reads tasks from the configured JSON file.
func (s *JSONStorage) LoadTasks() ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Task{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []models.Task{}, nil
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveTasks writes tasks to the configured JSON file.
func (s *JSONStorage) SaveTasks(tasks []models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
