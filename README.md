# cfg
Go package that defines a standard way to load configuration using viper.

Expects a dir `config` at the same level of the executable or a flag `--config` indicating a path to the configuration files

Expected default config files:
 - config/default.toml
 - config/local.toml

Suggested config files:
 - config/default.toml
 - config/local.toml
 - config/develop.toml
 - config/test.toml
 - config/staging.toml
 - config/production.toml

NOTE: Configuration files must be valid `toml` files.

`GOENV` env variable defines the configuration file that will be loaded


### Requirements
 - Viper

### Usage

```go
# main.go
package main

import (
    "fmt"

    "github.com/crabi/cfg"
    "github.com/spf13/viper"
)


func init() {
   cfg.Load()
}

func main() {
    fmt.Printf("%v", viper.GetString("SomeString"))
}
```

```bash
$ go build main
```
```bash
$ ls config
develop.toml local.toml production.toml staging.toml test.toml default.toml
```
```bash
$ GOENV=develop ./main
```
