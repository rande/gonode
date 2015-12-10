package core

import (
	"github.com/BurntSushi/toml"
	"github.com/rande/goapp"
)

func LoadConfiguration(path string, c interface{}) error {
	data, err := goapp.LoadConfigurationFromFile(path)

	goapp.PanicOnError(err)

	_, err = toml.Decode(data, c)

	goapp.PanicOnError(err)

	return nil
}
