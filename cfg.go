package cfg

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

const envModeVar = "GOENV"

var configPath string

// Load loads the configuration file depending on the Go environment mode
func Load() {
	flag.StringVar(&configPath, "config", "./config", "Configuration dir path")
	flag.Parse()

	viper.SetConfigType("toml")
	viper.SetConfigName("default")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("error reading default configuration file:", err)
	}

	configName, ok := os.LookupEnv(envModeVar)
	if !ok || configName == "" {
		configName = "local"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalln("error reading %s configuration file:", configName, err)
	}
}
