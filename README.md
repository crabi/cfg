# cfg
## About
Go package that defines a standard way to load configuration using viper.

## Disclaimer
This library does not standarizes how the required variables are loaded. A clear example of loading variables via
command line flags be found in `example/main.go`.

## Requirements
 - Viper

## Usage
Expected default config files:
 - config/default.yaml

Suggested config files:
 - config/local.yaml
 - config/develop.yaml
 - config/test.yaml
 - config/staging.yaml
 - config/production.yaml

NOTE: Configuration files must be valid `yaml` files.

For an example, see the directory `example/`.
