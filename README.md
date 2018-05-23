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
 - config/local.yaml

Suggested additional config files:
 - config/develop.yaml
 - config/test.yaml
 - config/staging.yaml
 - config/production.yaml

For an example, see the directory `example/`.
