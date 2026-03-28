package main

import (
	"fmt"
	"log"
)

// User Service Entry Point
//
// In a real system this would start a gRPC server. For this project,
// the service is wired in-process — the API gateway creates the service
// implementation directly.
//
// This main.go exists so the user-service package can be run independently
// for testing.
//
// Steps:
//  1. Create the in-memory store
//  2. Create the service (depends on store and JWT secret)
//  3. Optionally start a simple health-check HTTP server
//  4. Block until shutdown signal

func main() {
	// TODO: Create store, service, and optional health server
	fmt.Println("User Service")
	log.Println("User service is not yet implemented — use via API gateway")
}
