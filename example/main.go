package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/crabi/cfg"
)

func main() {
	// Retrieving the required values (config path and config file type) is up to the user
	var configPath string
	flag.StringVar(&configPath, "configPath", "", "Path to the configuration directory.")
	flag.Parse()

	// Config should be loaded before any other code functionality
	config, err := cfg.Load(&cfg.Params{Path: configPath})
	if err != nil {
		log.Fatal(err)
	}

	// Nested configurations can be retrieved using dot notation
	mainServiceConfig := config.Get("services.main")

	// Only to pretty print
	mainServiceConfigPretty, err := json.MarshalIndent(mainServiceConfig, "", "  ")
	fmt.Printf("%s", string(mainServiceConfigPretty))
}
