package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Bookstore API Server
//
// Steps to build:
//  1. Load configuration (from environment or defaults)
//  2. Create repositories (in-memory implementations)
//  3. Create services (business logic, depends on repositories)
//  4. Create handlers (HTTP layer, depends on services)
//  5. Set up the router and register routes
//  6. Wrap with middleware (logging, recovery, request ID, CORS)
//  7. Start the server with graceful shutdown
//
// Graceful shutdown pattern:
//   - Listen for SIGINT/SIGTERM in a goroutine
//   - Call server.Shutdown(ctx) when signal received
//   - Use a context with timeout so shutdown doesn't hang forever

func main() {
	// TODO: Load config
	// TODO: Create repositories
	// TODO: Create services
	// TODO: Create handlers
	// TODO: Set up routes

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// TODO: Register resource routes
	// TODO: Wrap mux with middleware

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		// TODO: Use context with timeout for shutdown
		if err := server.Close(); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
