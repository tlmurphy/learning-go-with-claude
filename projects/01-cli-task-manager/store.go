package main

// TaskStore manages a collection of tasks and persists them to a JSON file.
//
// TODO: Implement these methods:
//   - Load() error              — read tasks from the JSON file
//   - Save() error              — write tasks to the JSON file
//   - Add(title string, priority Priority) Task
//   - List(statusFilter *Status, priorityFilter *Priority) []Task
//   - Complete(id int) error
//   - Delete(id int) error

type TaskStore struct {
	FilePath string
	Tasks    []Task
	nextID   int
}

// NewTaskStore creates a new TaskStore with the given file path.
func NewTaskStore(filePath string) *TaskStore {
	return &TaskStore{
		FilePath: filePath,
		Tasks:    make([]Task, 0),
		nextID:   1,
	}
}
