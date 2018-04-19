package cfg

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	envModeVar      = "GOENV"
	localConfigName = "local"
)

var configPath string

func loadEnvSpecificConfig() {
	configName, ok := os.LookupEnv(envModeVar)
	if !ok || configName == "" {
		configName = localConfigName
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	if err := viper.MergeInConfig(); err != nil {
		if configName == localConfigName {
			log.Println(localConfigName, "file not found. Using only default config")
			return
		}

		log.Fatalln("error reading %s configuration file: %s", configName, err)
	}
}

// Load loads the configuration file depending on the Go environment mode
func Load(serviceName string) {
	flag.StringVar(&configPath, "config", "./config", "Configuration dir path")
	flag.Parse()

	viper.SetEnvPrefix(serviceName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("default")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("error reading default configuration file:", err)
	}

	loadEnvSpecificConfig()
}
