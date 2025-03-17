package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	ServerPort    int    `json:"server_port"`
	DatabasePath  string `json:"database_path"`
	BaseURL       string `json:"base_url"`
	MaxURLLength  int    `json:"max_url_length"`
	EnableLogging bool   `json:"enable_logging"`
}

// LoadSettings reads settings from a JSON file
func LoadSettings(filename string) (*Settings, error) {
	// First load from file if it exists
	var settings Settings

	fmt.Println("Loading settings from: ", filename)

	data, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			return nil, err
		}
	} else {
		fmt.Println("No settings file found, using defaults")
		settings = *GetDefaultSettings()
	}

	// Override with environment variables if they exist
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			settings.ServerPort = p
		}
	}

	if dbPath := os.Getenv("DATABASE_FILE"); dbPath != "" {
		settings.DatabasePath = dbPath
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		settings.BaseURL = baseURL
	}

	if maxLen := os.Getenv("MAX_URL_LENGTH"); maxLen != "" {
		if l, err := strconv.Atoi(maxLen); err == nil {
			settings.MaxURLLength = l
		}
	}

	if logging := os.Getenv("ENABLE_LOGGING"); logging != "" {
		settings.EnableLogging = strings.ToLower(logging) == "true"
	}

	return &settings, nil
}

// SaveSettings writes settings to a JSON file
func (s *Settings) SaveSettings(filename string) error {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// GetDefaultSettings returns default configuration values
func GetDefaultSettings() *Settings {
	return &Settings{
		ServerPort:    8080,
		DatabasePath:  "./urls.db",
		BaseURL:       "http://localhost:8080",
		MaxURLLength:  2048,
		EnableLogging: true,
	}
}
