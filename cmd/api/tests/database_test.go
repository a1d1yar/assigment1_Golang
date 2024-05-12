package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Define a struct to hold your test suite
type DatabaseTestSuite struct {
	db *sql.DB
}

// Setup function to initialize the test suite
func (suite *DatabaseTestSuite) Setup(t *testing.T) {
	// Open a connection to the test database
	db, err := sql.Open("postgres", "dbname=a.maratovDB user=postgres password=Aldiyar2004 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Assign the database connection to the test suite's db field
	suite.db = db
}

// Teardown function to clean up resources after the tests
func (suite *DatabaseTestSuite) Teardown(t *testing.T) {
	// Close the database connection
	err := suite.db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

// Test function to test database operations
func (suite *DatabaseTestSuite) TestDatabaseOperation(t *testing.T) {

	_, err := suite.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", "Test User", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

// Create a new instance of the test suite
var testSuite DatabaseTestSuite

// TestMain function to set up and tear down the test suite
func TestMain(m *testing.M) {
	// Set up the test suite
	testSuite.Setup(&testing.T{})

	// Run the tests

	testSuite.TestDatabaseOperation(&testing.T{})
	exitCode := m.Run()

	// Tear down the test suite
	testSuite.Teardown(&testing.T{})

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
