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

var (
	configPath string
	configType string
)

func loadEnvConfig(env string) error {
	viper.SetConfigName(env)
	viper.AddConfigPath(configPath)
	return viper.MergeInConfig()
}

// Load loads the configuration file depending on the Go environment mode
func Load(serviceName string) {
	flag.StringVar(&configPath, "config", "./config", "Configuration dir path")
	flag.StringVar(&configType, "configtype", "yaml", "Configuration format to use in files")
	flag.Parse()

	viper.SetConfigType(configType)

	viper.SetEnvPrefix(serviceName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("default")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("error reading default configuration file:", err)
	}

	envConfig, ok := os.LookupEnv(envModeVar)
	if !ok || envConfig == "" {
		envConfig = localConfigName
	}

	if err := loadEnvConfig(envConfig); err != nil {
		if envConfig == localConfigName {
			log.Println(localConfigName, "file not found. Using only default config")
			return
		}

		log.Fatalln("error reading %s configuration file: %s", envConfig, err)
	}
}
