package handler

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestHandleCommandNotFound(t *testing.T) {
    requestBody := `{"command":"command-that-does-not-exist"}`
    req, err := http.NewRequest(http.MethodPost, "/api/cmd", strings.NewReader(requestBody))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleCommand)

    handler.ServeHTTP(rr, req)

	// should be 404
    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code for a command not found: got %v want %v", status, http.StatusNotFound)
    }
}

func TestHandleCommandNonPostMethod(t *testing.T) {
    req, err := http.NewRequest(http.MethodGet, "/api/cmd", nil)
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleCommand)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("handler returned wrong status code for non-POST method: got %v want %v", status, http.StatusMethodNotAllowed)
    }
}