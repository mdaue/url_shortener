package main

import (
	"flag"
	"fmt"
	"os"
)

var config *Settings

func main() {
	var settingsFile string
	if env_filename := os.Getenv("SETTINGS_FILE"); env_filename != "" {
		settingsFile = env_filename
	} else {
		settingsFile = *flag.String("config", "urls.json", "path to settings file")
		flag.Parse()
	}

	var err error
	config, err = LoadSettings(settingsFile)
	if err != nil {
		fmt.Println("Error loading settings, using defaults")
		config = GetDefaultSettings()
		if err := config.SaveSettings(settingsFile); err != nil {
			fmt.Println("Error saving default settings:", err)
		}
	}

	Serve()
}
