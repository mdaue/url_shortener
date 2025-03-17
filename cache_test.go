package main

import (
	"testing"
)

func TestCreateCache(t *testing.T) {
	cache, err := createCache(10)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	if len(*cache.cache) != 0 {
		t.Errorf("Expected empty cache, got size %d", len(*cache.cache))
	}
}

func TestCacheURL(t *testing.T) {
	cache, _ := createCache(10)
	shortURL := "abc123"
	longURL := "https://example.com"

	cache.cacheURL(shortURL, longURL)

	if (*cache.cache)[shortURL] != longURL {
		t.Errorf("Expected URL %s, got %s", longURL, (*cache.cache)[shortURL])
	}
}

func TestGetURL(t *testing.T) {
	cache, _ := createCache(10)
	shortURL := "abc123"
	longURL := "https://example.com"

	cache.cacheURL(shortURL, longURL)

	got, err := cache.getURL(shortURL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got != longURL {
		t.Errorf("Expected URL %s, got %s", longURL, got)
	}

	// Test non-existent URL
	_, err = cache.getURL("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent URL")
	}
}

func TestPruneCache(t *testing.T) {
	cache, _ := createCache(10)

	// Add test data
	testData := map[string]string{
		"abc1": "url1",
		"abc2": "url2",
		"abc3": "url3",
		"abc4": "url4",
	}

	for k, v := range testData {
		cache.cacheURL(k, v)
	}

	initialSize := len(*cache.cache)
	cache.pruneCache()
	finalSize := len(*cache.cache)

	expectedSize := initialSize / 2
	if finalSize != expectedSize {
		t.Errorf("Expected size %d after pruning, got %d", expectedSize, finalSize)
	}
}
