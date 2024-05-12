package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEditUserInfoHandler(t *testing.T) {
	// Create a new instance of the application
	app := &application{}

	// Create a new request with a JSON body containing the user data to update
	reqBody := `{"name": "Updated User", "email": "updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/users/123", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the editUserInfoHandler with the mock request and recorder
	app.editUserInfoHandler(rr, req)

	// Check if the status code is OK (200)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// You can add more assertions here to check if the user data was updated as expected
}
