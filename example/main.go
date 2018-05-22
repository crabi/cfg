package main

import (
	"flag"
	"fmt"

	"github.com/crabi/cfg"
)

var (
	configPath string
	configName string
)

func init() {
	flag.StringVar(&configPath, "configPath", "./config", "Path to the configuration directory.")
	flag.StringVar(&configName, "configName", "local", "Running configuration name.")
}

func main() {
	flag.Parse()
	config, err := cfg.LoadConfig(configPath, configName)
	if err != nil {
		fmt.Printf("Error loading config: %+v", err)
		return
	}

	fmt.Printf("%+v", config)
}
