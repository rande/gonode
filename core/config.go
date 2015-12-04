package core

import (
	"github.com/BurntSushi/toml"
	"github.com/rande/goapp"
)

type ClientConfig struct {
	Server string `toml:"server"`
	Bind   string `toml:"bind"`
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{}
}

func LoadConfiguration(path string, c interface{}) error {
	data, err := goapp.LoadConfigurationFromFile(path)

	goapp.PanicOnError(err)

	_, err = toml.Decode(data, c)

	goapp.PanicOnError(err)

	return nil
}
