package main

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Environment app environment
type Environment string

// app environment
const (
	EnvProd Environment = "production"
	EnvStg              = "staging"
	EnvDev              = "development"
	EnvTest             = "test"
)

// AppConfig application config
type AppConfig struct {
	ServerPort  int    `toml:"server_port"`
	Environment string `toml:"environment"`
	Debug       bool   `toml:"debug"`
}

// NewAppConfig reads config file, and creates AppConfig
func NewAppConfig(path string) (*AppConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open config file: %s", path)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	var config AppConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		return nil, errors.Wrap(err, "failed to create AppConfig from file")
	}
	return &config, nil
}
