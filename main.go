package main

import (
	"github.com/joho/godotenv"
	"github.com/yuktea/golang-d/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Directly register the handler with the default serve mux.
	http.HandleFunc("/api/cmd", handler.HandleCommand)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get PORT from environment variables, use a default if not found
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	server := &http.Server{
		Addr: ":" + port, // Use the port from .env or default
	}

	// Starting the server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Server is starting on Port", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

    // Setup channel to listen for interrupt or termination signals from the OS
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    // Block until a signal is received ie, quit has value
    <-quit
    log.Println("Server is shutting down...")

    // Attempt to gracefully shut down the server
    if err := server.Shutdown(nil); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited gracefully")
}
