package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Page holds the data to be passed to HTML templates
type Page struct {
	Title       string
	URLs        []URL // This should match the type you're using in your DB queries
	CurrentTime string
}

// HomeHandler handles the root path and serves the main page
type HomeHandler struct {
	db    *sql.DB
	cache *Cache
}

// ServeHTTP implements the http.Handler interface
func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the 10 most recent URLs from the database
	urls, err := queryRecentURLs(h.db, 10)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	// Create page data
	page := Page{
		Title:       "URL Shortener",
		URLs:        urls,
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/index.html", "templates/url_list.html", "templates/url_row.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// URLFormHandler handles the form submission for creating new short URLs
type URLFormHandler struct {
	db    *sql.DB
	cache *Cache
}

// ServeHTTP implements the http.Handler interface
func (h URLFormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract form data
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	/*
		Check if the URL is valid
	*/
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	if len(originalURL) > config.MaxURLLength {
		http.Error(w, "URL exceeds maximum length", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	blockedPatterns := []string{".exe", "javascript:", "data:"}
	for _, pattern := range blockedPatterns {
		if strings.Contains(strings.ToLower(originalURL), pattern) {
			http.Error(w, "URL contains forbidden content", http.StatusBadRequest)
			return
		}
	}

	// Create short URL
	shortUrl, err := shorten(originalURL)
	if err != nil {
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	_, err = createURL(h.db, originalURL, shortUrl, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	// Add to cache
	h.cache.cacheURL(shortUrl, originalURL)

	// Return just the new URL row as HTML for HTMX to insert
	urlData := URL{
		Name:          originalURL,
		CreatedAt:     time.Now(),
		Short:         shortUrl,
		RequestedFrom: r.RemoteAddr,
		Clicks:        0,
	}

	// Parse and execute the partial template
	tmpl, err := template.ParseFiles("templates/url_row.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set HTMX specific headers
	w.Header().Set("HX-Trigger", "urlAdded")
	err = tmpl.ExecuteTemplate(w, "url_row", urlData)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// RefreshHandler handles refreshing the URL list
type RefreshHandler struct {
	db *sql.DB
}

// ServeHTTP implements the http.Handler interface
func (h RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the 10 most recent URLs from the database
	urls, err := queryRecentURLs(h.db, 10)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	// Create page data
	page := struct {
		URLs        []URL
		CurrentTime string
	}{
		URLs:        urls,
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/url_list.html", "templates/url_row.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "url_list", page)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// StaticFileHandler serves static files (CSS, JS, etc.)
func StaticFileHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}

type QueryHandler struct {
	db    *sql.DB
	cache *Cache
}

func (qh QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Extract short URL from path
		shortURL := strings.TrimPrefix(r.URL.Path, "/q/")
		if shortURL != "" {
			url, err := qh.cache.getURL(shortURL)
			if err != nil {
				url, err := queryShortURL(qh.db, shortURL)
				if err != nil {
					http.NotFound(w, r)
					return
				}
				qh.cache.cacheURL(shortURL, url.Name)
			}
			// Cache the URL after successfully retrieving it
			qh.cache.cacheURL(shortURL, url)
			// Redirect to long URL for all HTTP methods
			addClicks(qh.db, shortURL)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// SetupRoutes sets up the routes for the web application
func SetupRoutes(db *sql.DB, cache *Cache) {
	// Add new handlers for the web frontend
	http.Handle("/", HomeHandler{db: db, cache: cache})
	http.Handle("/create", URLFormHandler{db: db, cache: cache})
	http.Handle("/refresh", RefreshHandler{db: db})
	http.Handle("/static/", StaticFileHandler())

	// Add the existing REST API
	http.Handle("/s/", URLFormHandler{db: db, cache: cache})
	http.Handle("/q/", QueryHandler{db: db, cache: cache})
}

// Serve sets up and starts the server
func Serve() {
	db, _ := openDatabase()
	cache, _ := createCache(1024)

	SetupRoutes(db, cache)

	println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
