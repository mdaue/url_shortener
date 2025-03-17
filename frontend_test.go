package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db
}

func TestHomeHandler(t *testing.T) {
	db := setupTestDB(t)
	cache, _ := createCache(10)
	handler := HomeHandler{db: db, cache: cache}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "URL Shortener") {
		t.Error("Homepage should contain 'URL Shortener'")
	}
}

func TestURLFormHandler(t *testing.T) {
	db := setupTestDB(t)
	cache, _ := createCache(10)
	handler := URLFormHandler{db: db, cache: cache}

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid URL",
			method:         "POST",
			url:            "https://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         "GET",
			url:            "https://example.com",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Empty URL",
			method:         "POST",
			url:            "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Malicious URL",
			method:         "POST",
			url:            "javascript:alert(1)",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("url", tt.url)
			req := httptest.NewRequest(tt.method, "/create", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	db := setupTestDB(t)
	handler := RefreshHandler{db: db}

	req := httptest.NewRequest("GET", "/refresh", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStaticFileHandler(t *testing.T) {
	handler := StaticFileHandler()

	req := httptest.NewRequest("GET", "/static/css/styles.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Error("Static file handler should serve files from static directory")
	}
}

func TestSetupRoutes(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	staticDir := filepath.Join(testDir, "static")
	os.MkdirAll(staticDir, 0755)
	os.WriteFile(filepath.Join(staticDir, "style.css"), []byte("body { color: black; }"), 0644)

	db, err := openDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	cache, _ := createCache(10)

	// Create test URL
	shortURL := "testshort"
	originalURL := "https://example.com"
	_, err = createURL(db, originalURL, shortURL, "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	cache.cacheURL(shortURL, originalURL)

	mux := http.NewServeMux()
	mux.Handle("/", HomeHandler{db: db, cache: cache})
	mux.Handle("/create", URLFormHandler{db: db, cache: cache})
	mux.Handle("/refresh", RefreshHandler{db: db})
	mux.Handle("/static/", http.FileServer(http.Dir(testDir)))
	mux.Handle("/s/", URLFormHandler{db: db, cache: cache})
	mux.Handle("/q/", QueryHandler{db: db, cache: cache})

	tests := []struct {
		path   string
		method string
		body   string
		want   int
	}{
		{"/", "GET", "", http.StatusOK},
		{"/create", "POST", "url=https://example.com", http.StatusOK},
		{"/refresh", "GET", "", http.StatusOK},
		{"/static/style.css", "GET", "", http.StatusOK},
		{"/s", "POST", "url=https://example.com", http.StatusMovedPermanently},
		{"/q/testshort", "GET", "", http.StatusMovedPermanently},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.method == "POST" {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("got %d, want %d for %s", w.Code, tt.want, tt.path)
			}
		})
	}
}
