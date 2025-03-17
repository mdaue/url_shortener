package main

import (
	"encoding/json"
	"os"
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
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
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
