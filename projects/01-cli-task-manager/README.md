# Project 01: CLI Task Manager

**Prerequisite modules:** 01-08 (Foundations)

## Overview

Build a command-line task manager that stores tasks in a JSON file. This project
exercises everything from the Foundations modules: variables, control flow,
functions, collections, structs, interfaces, pointers, and error handling.

When complete you will have a useful tool you can actually run day-to-day.

## Requirements

### Core Features
- **Add** tasks with a title and optional priority (low, medium, high)
- **List** all tasks, with filtering by status and/or priority
- **Complete** a task (mark it done)
- **Delete** a task
- **Persist** tasks to a JSON file and reload on startup

### Technical Requirements
- Define a `Task` struct with proper types
  - Use `iota`-based constants for Priority (Low, Medium, High) and Status (Pending, Done)
  - Include created-at and completed-at timestamps
- Implement the `fmt.Stringer` interface on `Task` so tasks display nicely
- Parse command-line arguments with `os.Args` or the `flag` package
- Use JSON struct tags for clean serialization (`encoding/json`)
- Handle all file I/O errors gracefully (missing file on first run, permission errors, corrupt JSON, etc.)

### Example Usage

```
$ taskmanager add "Write unit tests" --priority high
Task added: #1 Write unit tests [high]

$ taskmanager add "Buy groceries"
Task added: #2 Buy groceries [medium]

$ taskmanager list
#1 [ ] Write unit tests      (high)   created 2026-03-27
#2 [ ] Buy groceries         (medium) created 2026-03-27

$ taskmanager complete 1
Completed: #1 Write unit tests

$ taskmanager list --status done
#1 [x] Write unit tests      (high)   completed 2026-03-27

$ taskmanager delete 2
Deleted: #2 Buy groceries
```

## Hints

<details>
<summary>Architecture hint</summary>

Consider a `TaskStore` struct that holds a `[]Task` and the file path.
Attach methods to `TaskStore` for Add, List, Complete, Delete, Save, and Load.
This keeps your `main()` clean — it just parses arguments and calls the right
method on the store.

</details>

<details>
<summary>Serialization hint</summary>

The `encoding/json` package can marshal/unmarshal your `[]Task` directly.
Use `json.MarshalIndent` for human-readable output.
Make sure your struct fields are exported (capitalized) and have `json:"..."` tags.

</details>

<details>
<summary>Timestamps hint</summary>

`time.Now()` gives you the current time. Store it in your Task struct and
format it with `time.Format("2006-01-02")` for display. The completed-at
field can be a `*time.Time` (pointer) so it is nil when the task is not done.

</details>

<details>
<summary>ID generation hint</summary>

A simple approach: when loading from file, find the max ID across all tasks
and store it as `nextID` in your `TaskStore`. Increment it each time you add
a task.

</details>

## Stretch Goals

- **Due dates** — add an optional due date and highlight overdue tasks in the list output
- **Multiple output formats** — support `--format table`, `--format json`, and `--format csv`
- **Categories / tags** — let users assign tags and filter by them
- **Undo last action** — keep a simple history and support `taskmanager undo`
