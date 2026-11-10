package models

import "time"

// Priority defines the priority level of a task.
type Priority string

const (
	PriorityLow    Priority = "Low"
	PriorityMedium Priority = "Medium"
	PriorityHigh   Priority = "High"
)

// Status defines the current progress state of a task.
type Status string

const (
	StatusPending    Status = "Pending"
	StatusInProgress Status = "In Progress"
	StatusCompleted  Status = "Completed"
)

// Task represents a task entity within TaskFlow.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	Category    string    `json:"category"`
	DueDate     string    `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
