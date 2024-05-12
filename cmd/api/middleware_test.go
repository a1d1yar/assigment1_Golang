package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverPanicMiddleware(t *testing.T) {
	// Define a mock handler function for testing purposes
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Create an instance of the application
	app := &application{}

	// Call the recoverPanic middleware with the mock handler
	recoveredHandler := app.recoverPanic(handler)
	recoveredHandler.ServeHTTP(rr, req)

	// Check if the status code is OK (200)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// You can add more assertions here to check if the middleware is working as expected
}

func TestRateLimitMiddleware(t *testing.T) {
	// Define a mock handler function for testing purposes
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Create an instance of the application
	app := &application{}

	// Call the rateLimit middleware with the mock handler
	rateLimitedHandler := app.rateLimit(handler)
	rateLimitedHandler.ServeHTTP(rr, req)

	// Check if the status code is OK (200)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// You can add more assertions here to check if the middleware is working as expected
}
