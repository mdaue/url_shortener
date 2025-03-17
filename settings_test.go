package main

import (
	"os"
	"testing"
)

func TestSettingsLoadAndSave(t *testing.T) {
	// Test settings
	testSettings := &Settings{
		ServerPort:    9090,
		DatabasePath:  "test.db",
		BaseURL:       "http://test.com",
		MaxURLLength:  1024,
		EnableLogging: true,
	}

	// Test file name
	testFile := "test_settings.json"

	// Save settings
	err := testSettings.SaveSettings(testFile)
	if err != nil {
		t.Fatalf("Failed to save settings: %v", err)
	}

	// Load settings
	loadedSettings, err := LoadSettings(testFile)
	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}

	// Compare values
	if loadedSettings.ServerPort != testSettings.ServerPort {
		t.Errorf("ServerPort mismatch: got %d, want %d", loadedSettings.ServerPort, testSettings.ServerPort)
	}
	if loadedSettings.DatabasePath != testSettings.DatabasePath {
		t.Errorf("DatabasePath mismatch: got %s, want %s", loadedSettings.DatabasePath, testSettings.DatabasePath)
	}
	if loadedSettings.BaseURL != testSettings.BaseURL {
		t.Errorf("BaseURL mismatch: got %s, want %s", loadedSettings.BaseURL, testSettings.BaseURL)
	}
	if loadedSettings.MaxURLLength != testSettings.MaxURLLength {
		t.Errorf("MaxURLLength mismatch: got %d, want %d", loadedSettings.MaxURLLength, testSettings.MaxURLLength)
	}
	if loadedSettings.EnableLogging != testSettings.EnableLogging {
		t.Errorf("EnableLogging mismatch: got %v, want %v", loadedSettings.EnableLogging, testSettings.EnableLogging)
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestGetDefaultSettings(t *testing.T) {
	defaults := GetDefaultSettings()

	expectedPort := 8080
	if defaults.ServerPort != expectedPort {
		t.Errorf("Default ServerPort mismatch: got %d, want %d", defaults.ServerPort, expectedPort)
	}

	expectedDBPath := "./urls.db"
	if defaults.DatabasePath != expectedDBPath {
		t.Errorf("Default DatabasePath mismatch: got %s, want %s", defaults.DatabasePath, expectedDBPath)
	}

	expectedBaseURL := "http://localhost:8080"
	if defaults.BaseURL != expectedBaseURL {
		t.Errorf("Default BaseURL mismatch: got %s, want %s", defaults.BaseURL, expectedBaseURL)
	}

	expectedMaxLength := 2048
	if defaults.MaxURLLength != expectedMaxLength {
		t.Errorf("Default MaxURLLength mismatch: got %d, want %d", defaults.MaxURLLength, expectedMaxLength)
	}

	expectedLogging := true
	if defaults.EnableLogging != expectedLogging {
		t.Errorf("Default EnableLogging mismatch: got %v, want %v", defaults.EnableLogging, expectedLogging)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := LoadSettings("nonexistent.json")
	if err != nil {
		t.Error("Unexpected error when loading non-existent file ", err)
	}
}
