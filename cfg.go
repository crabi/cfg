package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configType        = "yaml"
	defaultConfigName = "default"
)

// LoadConfig loads a config directory, the default config file and overrides with a given configuration.
// It returns a generic configuration, which ideally will be parsed into a struct.
func LoadConfig(configPath, configName string) (map[string]interface{}, error) {
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	// Load default config
	viper.SetConfigName(defaultConfigName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading \"%s\" configuration", defaultConfigName)
	}

	// Override with specific config
	if configName != "" {
		viper.SetConfigName(configName)
		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error overriding \"%s\" configuration", configName)
		}
	}

	return viper.AllSettings(), nil
}
