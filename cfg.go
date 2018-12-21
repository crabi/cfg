package cfg

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	envNameKey = "GOENV"

	defaultConfigPath = "./config"
	defaultFileType   = "yaml"

	defaultConfigEnv = "default"
	localConfigEnv   = "local"
	envVarConfigFile = "environment-variables"

	dotEnvFile = ".env"
)

// Params is used to load a new environment
type Params struct {
	Path          string
	FileType      string
	RequireDotEnv bool
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

// Config gives access to the loaded configuration
type Config interface {
	Get(string) interface{}
	AllSettings() map[string]interface{}
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

func (c *configWrapper) loadEnvConfigFile() {
	tempV := viper.New()

	tempV.SetConfigType(c.params.FileType)
	tempV.SetConfigName(envVarConfigFile)
	tempV.AddConfigPath(c.params.Path)
	if err := tempV.ReadInConfig(); err != nil {
		return
	}
	keys := tempV.AllKeys()
	for _, key := range keys {
		envVarName := tempV.GetString(key)
		if value, exist := os.LookupEnv(envVarName); exist {
			if isBool(value) {
				bValue, _ := strconv.ParseBool(value)
				c.v.Set(key, bValue)
			} else if isInt(value) {
				iValue, _ := strconv.Atoi(value)
				c.v.Set(key, iValue)
			} else if isFloat(value) {
				fValue, _ := strconv.ParseFloat(value, 64)
				c.v.Set(key, fValue)
			} else {
				c.v.Set(key, value)
			}
		}
	}

	return
}

// Get returns a map[string]interface{} of a given key
func (c *configWrapper) Get(key string) interface{} {
	d := c.v.AllSettings()
	v := interface{}(d)
	path := strings.Split(key, ".")
	for _, key := range path {
		switch v.(type) {
		case map[string]interface{}:
			v = v.(map[string]interface{})[key]
		default:
			return nil
		}

	}
	return v
}

func isInt(data string) bool {
	if _, err := strconv.Atoi(data); err == nil {
		return true
	}
	return false
}

func isFloat(data string) bool {
	if _, err := strconv.ParseFloat(data, 64); err == nil {
		return true
	}
	return false
}

func isBool(data string) bool {
	if _, err := strconv.ParseBool(data); err == nil {
		return true
	}
	return false
}

func (c *configWrapper) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

func loadDotEnv(required bool) error {
	_, err := os.Stat(dotEnvFile)
	dotEnvExist := !os.IsNotExist(err)

	if required && !dotEnvExist {
		return errors.New("missing required " + dotEnvFile)
	}

	if !required && !dotEnvExist {
		return nil
	}

	return godotenv.Load(dotEnvFile)
}

// Load loads a config directory, the default config file and overrides with a given configuration.
// It returns a generic configuration, which ideally will be parsed into a struct.
func Load(params *Params) (Config, error) {
	if params == nil {
		return nil, errors.New("nil params")
	}

	if err := loadDotEnv(params.RequireDotEnv); err != nil {
		return nil, err
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
	if err := c.mergeEnv(envName); err != nil && envName != localConfigEnv {
		return nil, fmt.Errorf("error merging \"%s\" configuration file", envName)
	}
	c.loadEnvConfigFile()

	return c, nil
}
