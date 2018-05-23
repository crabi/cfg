package cfg

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	envModeVar        = "GOENV"
	defaultFileType   = "yaml"
	defaultConfigEnv  = "default"
	defaultConfigPath = "./config"
	localConfigEnv    = "local"
)

type Cfg struct {
	SvcName  string
	Path     string
	FileType string
}

type cfg interface {
	getEnv() string
	loadEnv(string) error
	mergeEnv(string) error
	setCoreUp()
	setDefaults()
}

// Load loads a config directory, the default config file and overrides with a given configuration.
// It returns a generic configuration, which ideally will be parsed into a struct.
func Load(c cfg) {
	c.setDefaults()
	c.setCoreUp()

	if err := c.loadEnv(defaultConfigEnv); err != nil {
		log.Fatalln("error reading default configuration file:", err)
	}

	env := c.getEnv()
	if err := c.mergeEnv(env); err != nil {
		log.Fatalln("error reading %s configuration file: %s", env, err)
	}
}

func (c *Cfg) loadEnv(env string) error {
	viper.SetConfigName(env)
	viper.AddConfigPath(c.Path)
	return viper.ReadInConfig()
}

func (c *Cfg) mergeEnv(env string) error {
	viper.SetConfigName(env)
	viper.AddConfigPath(c.Path)
	return viper.MergeInConfig()
}

func (c *Cfg) setDefaults() {
	if c.Path == "" {
		c.Path = defaultConfigPath
	}

	if c.FileType == "" {
		c.FileType = defaultFileType
	}
}

func (c *Cfg) setCoreUp() {
	viper.SetConfigType(c.FileType)
	viper.SetEnvPrefix(c.SvcName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func (c *Cfg) getEnv() string {
	env, ok := os.LookupEnv(envModeVar)
	if !ok || env == "" {
		env = localConfigEnv
	}
	return env
}

func Get(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
