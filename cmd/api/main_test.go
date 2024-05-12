package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func (app *application) TestMain(t *testing.T) {
	// Create a new request to the API server
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a new API server and handle the request
	handler := http.HandlerFunc(app.createModuleInfoHandler)
	handler.ServeHTTP(rr, req)

	// Check if the response status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body contains the expected message
	expected := `{"message":"Hello, world!"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
