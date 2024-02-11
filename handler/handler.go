package handler

import (
    "encoding/json"
    "net/http"
    "os/exec"
	"log"
	"bytes"
)

type RequestData struct {
    Command string
}

type CommandResponse struct {
    Stdout    string `json:"stdout"`
    Stderr    string `json:"stderr"` 
    ErrorCode int    `json:"errorCode"` 
}

const ExitCodeNotFound = 127

func HandleCommand(w http.ResponseWriter, r *http.Request) {

    // Check if the request method is not POST
	if checkIfNotPost(w, r) {
		return
	}

    var data RequestData
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        log.Printf("Error parsing JSON: %v", err)
        respondWithJSON(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    if data.Command == "" {
        respondWithJSON(w, http.StatusBadRequest, "Command not provided")
        return
    }

    cmd := exec.Command("sh", "-c", data.Command)
    var stdoutBuf, stderrBuf bytes.Buffer
    cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

    execErr := cmd.Run()
    response := CommandResponse{ // Assuming CommandResponse is defined to include Stdout, Stderr, and ErrorCode
        Stdout:    stdoutBuf.String(),
        Stderr:    stderrBuf.String(),
		ErrorCode: 0,
    }

    if execErr != nil {
        // Adjust ErrorCode based on the error type
        if exitErr, ok := execErr.(*exec.ExitError); ok {
            response.ErrorCode = exitErr.ExitCode()
        } else {
            response.ErrorCode = -1 // An unexpected error type
        }
        log.Printf("Error executing command: '%s', Exit Code: %d, Error: %v, Stderr: %s\n", data.Command, response.ErrorCode, execErr, stderrBuf.String())
        
		// dont get
        statusCode := http.StatusInternalServerError
        if response.ErrorCode == ExitCodeNotFound {
            statusCode = http.StatusNotFound
        }
        respondWithJSON(w, statusCode, response)
        return
    }

    // For successful execution:
    respondWithJSON(w, http.StatusOK, response)
}


func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}


func checkIfNotPost(w http.ResponseWriter, r *http.Request) bool {
    if r.Method != http.MethodPost {
        respondWithJSON(w, http.StatusMethodNotAllowed, "Method not allowed (can only use POST)")
        return true
    }
    return false
}


//test with
// curl -X POST 'http://localhost:5555/api/cmd' \                             
// -H 'Content-Type: application/json' \
// -d '{"command": "echo Hello World"}'
// curl -X POST 'http://localhost:5555/api/cmd?command=echo%20Hello%20World' 

