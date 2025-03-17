package main

import (
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("Database connection not active: %v", err)
	}
}

func TestCreateURL(t *testing.T) {
	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test data
	name := "https://example.com"
	short := "abc123"
	requestedFrom := "127.0.0.1"

	url, err := createURL(db, name, short, requestedFrom)
	if err != nil {
		t.Errorf("Failed to create URL: %v", err)
	}

	if url.Name != name {
		t.Errorf("Expected name %s, got %s", name, url.Name)
	}
	if url.Short != short {
		t.Errorf("Expected short %s, got %s", short, url.Short)
	}
	if url.RequestedFrom != requestedFrom {
		t.Errorf("Expected requested_from %s, got %s", requestedFrom, url.RequestedFrom)
	}
}

func TestQueryURLs(t *testing.T) {
	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	urls, err := queryURLs(db)
	if err != nil {
		t.Errorf("Failed to query URLs: %v", err)
	}

	// Verify we can retrieve URLs
	if len(urls) == 0 {
		t.Log("No URLs found in database")
	}
}

func TestQueryURLsFromRequested(t *testing.T) {
	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	requestedFrom := "127.0.0.1"
	urls, err := queryURLsFromRequested(db, requestedFrom)
	if err != nil {
		t.Errorf("Failed to query URLs by requested_from: %v", err)
	}

	// Verify all returned URLs have the correct requested_from value
	for _, url := range urls {
		if url.RequestedFrom != requestedFrom {
			t.Errorf("Expected requested_from %s, got %s", requestedFrom, url.RequestedFrom)
		}
	}
}

func TestQueryRecentURLs(t *testing.T) {
	db, _ := openDatabase()
	defer db.Close()

	// Insert test data with different timestamps
	testURLs := []struct {
		name          string
		short         string
		requestedFrom string
	}{
		{"https://first.com", "abc123", "127.0.0.1"},
		{"https://second.com", "def456", "127.0.0.1"},
		{"https://third.com", "ghi789", "127.0.0.1"},
	}

	for _, u := range testURLs {
		createURL(db, u.name, u.short, u.requestedFrom)
	}

	// Test with limit of 2
	limit := 2
	urls, err := queryRecentURLs(db, limit)

	if err != nil {
		t.Fatalf("queryRecentURLs failed: %v", err)
	}

	if len(urls) != limit {
		t.Errorf("Expected %d URLs, got %d", limit, len(urls))
	}

	// Verify order - should be newest first
	for i := 1; i < len(urls); i++ {
		if urls[i-1].CreatedAt.Before(urls[i].CreatedAt) {
			t.Errorf("URLs not in correct order. Expected descending by creation time")
		}
	}
}
