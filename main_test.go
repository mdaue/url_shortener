package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Setup default settings
	config = GetDefaultSettings()

	// Setup test database
	db, err := openDatabase()
	if err != nil {
		os.Exit(1)
	}

	// Create test table
	_, err = db.Exec(`
		DROP TABLE IF EXISTS urls;
        CREATE TABLE IF NOT EXISTS urls (
            name TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            short TEXT NOT NULL,
            requested_from TEXT NOT NULL,
			clicks INTEGER DEFAULT 0
        )
    `)
	if err != nil {
		db.Close()
		os.Exit(1)
	}
	db.Close()

	// Start server in background
	go func() {
		Serve()
	}()

	// Wait for server to be ready
	client := http.Client{Timeout: time.Second}
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := client.Get("http://localhost:8080/u")
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Run all tests
	code := m.Run()

	// Cleanup
	os.Remove("urls.sql")
	os.Exit(code)
}
