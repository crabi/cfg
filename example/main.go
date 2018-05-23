package main

import (
	"flag"
	"fmt"

	"github.com/crabi/cfg"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "Path to the configuration directory.")
	cfg.Load(&cfg.Cfg{
		SvcName: "main",
		Path:    configPath,
	})
}

func main() {
	svcCfg := cfg.Get("services.main")
	fmt.Printf("%+v", svcCfg)
}
