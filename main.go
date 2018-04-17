package cfg

import (
	"os"
	"log"
	"flag"

	"github.com/spf13/viper"
)

const EnvModeVar = "GOENV"

var configPath string

func Load() {
	flag.StringVar(&configPath, "config", "./config", "Configuration dir path")
	flag.Parse()

	viper.SetConfigType("toml")
	viper.SetConfigName("default")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("error reading default configuration file:", err)
	}

	configName, ok := os.LookupEnv(EnvModeVar)
	if !ok || configName == "" {
		configName = "local"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalln("error reading %s configuration file:", configName, err)
	}
}
