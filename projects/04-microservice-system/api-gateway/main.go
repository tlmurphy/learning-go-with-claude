package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// API Gateway Entry Point
//
// The gateway exposes a REST API and delegates to the UserService.
//
// Steps:
//  1. Create the UserService (or connect to it)
//  2. Create handlers (depends on UserService)
//  3. Set up routes with JWT middleware on protected endpoints
//  4. Start server with graceful shutdown

func main() {
	// TODO: Create user service (in-process for now)
	// TODO: Create handlers
	// TODO: Set up routes

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// TODO: Register API routes
	// POST /api/v1/register
	// POST /api/v1/login
	// GET  /api/v1/profile   (auth required)
	// PUT  /api/v1/profile   (auth required)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down gateway...")
		if err := server.Close(); err != nil {
			log.Printf("Gateway shutdown error: %v", err)
		}
	}()

	log.Printf("API Gateway listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gateway error: %v", err)
	}
}
