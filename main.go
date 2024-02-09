package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func handleCommand(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the command from query param or JSON body
	command := r.URL.Query().Get("command")
	if command == "" {
		http.Error(w, "Command not provided", http.StatusBadRequest)
		return
	}

	// Execute command
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		// Check if the command is not found
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 127 {
			http.Error(w, fmt.Sprintf("Command '%s' not found", command), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// output as response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func main() {
	http.HandleFunc("/api/cmd", handleCommand)
	log.Fatal(http.ListenAndServe(":5555", nil))
}
