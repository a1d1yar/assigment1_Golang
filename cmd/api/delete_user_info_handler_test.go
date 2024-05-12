package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUserInfoHandler(t *testing.T) {
	// Create a new instance of the application
	app := &application{}

	// Create a new request
	req := httptest.NewRequest(http.MethodDelete, "/v1/users/123", nil)

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the deleteUserInfoHandler with the mock request and recorder
	app.deleteUserInfoHandler(rr, req)

	// Check if the status code is OK (200)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// You can add more assertions here to check if the user was deleted as expected
}
