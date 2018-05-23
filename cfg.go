package cfg

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	envNameKey = "GOENV"

	defaultConfigPath = "./config"
	defaultConfigType = "viper"
	defaultFileType   = "yaml"

	defaultConfigEnv = "default"
	localConfigEnv   = "local"
)

// Params is used to load a new environment
type Params struct {
	Path     string
	FileType string
}

// SetDefaults sets the default values of a Config fields
func (c *Params) SetDefaults() {
	if c.Path == "" {
		c.Path = defaultConfigPath
	}
	if c.FileType == "" {
		c.FileType = defaultFileType
	}
}

type configWrapper struct {
	v      *viper.Viper
	params *Params
}

func (c *configWrapper) getEnv() string {
	envName := os.Getenv(envNameKey)
	if envName == "" {
		envName = localConfigEnv
	}
	return envName
}

func (c *configWrapper) loadEnv(envName string) error {
	c.v.SetConfigType(c.params.FileType)
	c.v.SetConfigName(envName)
	c.v.AddConfigPath(c.params.Path)
	return c.v.ReadInConfig()
}

func (c *configWrapper) mergeEnv(envName string) error {
	c.v.SetConfigType(c.params.FileType)
	c.v.SetConfigName(envName)
	c.v.AddConfigPath(c.params.Path)
	return c.v.MergeInConfig()
}

// Get returns a map[string]interface{} of a given key
func (c *configWrapper) Get(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

// Load loads a config directory, the default config file and overrides with a given configuration.
// It returns a generic configuration, which ideally will be parsed into a struct.
func Load(params *Params) (*configWrapper, error) {
	if params == nil {
		return nil, errors.New("nil params")
	}

	c := &configWrapper{
		v:      viper.New(),
		params: &(*params),
	}
	c.params.SetDefaults()

	if err := c.loadEnv(defaultConfigEnv); err != nil {
		return nil, fmt.Errorf("error reading \"default\" configuration file")
	}

	envName := c.getEnv()
	if err := c.mergeEnv(envName); err != nil {
		return nil, fmt.Errorf("error merging \"%s\" configuration file", envName)
	}

	return c, nil
}
