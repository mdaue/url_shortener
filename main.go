package main

import (
	"flag"
	"fmt"
)

var config *Settings

func main() {
	settingsFile := flag.String("config", "urls.json", "path to settings file")
	flag.Parse()

	var err error
	config, err = LoadSettings(*settingsFile)
	if err != nil {
		fmt.Println("Error loading settings, using defaults")
		config = GetDefaultSettings()
		if err := config.SaveSettings(*settingsFile); err != nil {
			fmt.Println("Error saving default settings:", err)
		}
	}

	Serve()
}
