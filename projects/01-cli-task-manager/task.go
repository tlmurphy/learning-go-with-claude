package main

import "time"

// Priority represents task priority levels.
// TODO: Define as iota-based constants (Low, Medium, High).
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

// Status represents task completion status.
// TODO: Define as iota-based constants (Pending, Done).
type Status int

const (
	StatusPending Status = iota
	StatusDone
)

// Task represents a single task in the task manager.
// TODO: Add JSON struct tags to all fields.
// TODO: Implement the fmt.Stringer interface.
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
