package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/a1d1yar/assingment1_Golang/internal/jsonlog"
)

func TestJSONLogger_PrintInfo(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "test_jsonlog_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Create a new JSONLogger instance
	logger := jsonlog.New(tmpfile, jsonlog.LevelInfo)

	// Log an INFO message
	logger.PrintInfo("This is an info message", nil)

	// Read the content of the temporary file
	content, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the log message was written correctly
	if !strings.Contains(string(content), `"level":"INFO"`) {
		t.Error("Expected INFO log level, got:", string(content))
	}
}

func TestJSONLogger_PrintError(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "test_jsonlog_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Create a new JSONLogger instance
	logger := jsonlog.New(tmpfile, jsonlog.LevelError)

	// Log an ERROR message
	logger.PrintError(errors.New("This is an error message"), nil)

	// Read the content of the temporary file
	content, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the log message was written correctly
	if !strings.Contains(string(content), `"level":"ERROR"`) {
		t.Error("Expected ERROR log level, got:", string(content))
	}
}
