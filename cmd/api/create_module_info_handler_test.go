package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func (app *application) TestCreateModuleInfoHandlerMissingField(t *testing.T) {

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/moduleinfoTEST", app.createModuleInfoHandler)

	// Create a new request with a JSON body missing the "name" field
	type payload struct {
		ModuleName     string `json:"module_name"`
		ModuleDuration int    `json:"module_duration"`
		ExamType       string `json:"exam_type"`
		Version        int    `json:"version"`
	}

	a := payload{
		ModuleName:     "This is test module",
		ModuleDuration: 6,
		ExamType:       "running",
		Version:        1,
	}

	expectedData, _ := json.Marshal(a)
	fmt.Println("Marshaled data: ", expectedData)
	req, err := http.NewRequest(http.MethodPost, "/v1/moduleinfoTEST", strings.NewReader(`{
		"module_name": "AITU2",
		"module_duration": 6,
		"exam_type": "running",
		"version": 1
	}`))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new ResponseRecorder, which satisfies the http.ResponseWriter interface
	rr := httptest.NewRecorder()

	// Call the handler function with the request and ResponseRecorder
	router.ServeHTTP(rr, req)

	// Assert that the response has a status code of http.StatusBadRequest
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// Assert that the response body contains an error message indicating the missing field
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	errorMessage := response["error"]
	if errorMessage != "missing required field: name" {
		t.Errorf("Expected error message to be 'missing required field: name', got %s", errorMessage)
	}
}
