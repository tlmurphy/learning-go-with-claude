package main

import (
	"fmt"
	"os"
)

// CLI Task Manager
//
// Your entry point. Parse command-line arguments and dispatch to the
// appropriate TaskStore method.
//
// Suggested commands:
//   taskmanager add "title" [--priority high|medium|low]
//   taskmanager list [--status pending|done] [--priority high|medium|low]
//   taskmanager complete <id>
//   taskmanager delete <id>
//
// Steps:
//  1. Create a TaskStore (with a file path like "tasks.json")
//  2. Load existing tasks from the file
//  3. Parse os.Args (or use the flag package) to determine the command
//  4. Execute the command
//  5. Save tasks back to the file if anything changed

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Create a TaskStore, load tasks, and dispatch commands.
	fmt.Println("Task Manager - not yet implemented")
}

func printUsage() {
	fmt.Println("Usage: taskmanager <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <title> [--priority low|medium|high]   Add a new task")
	fmt.Println("  list [--status pending|done] [--priority]  List tasks")
	fmt.Println("  complete <id>                              Mark a task as done")
	fmt.Println("  delete <id>                                Delete a task")
}
