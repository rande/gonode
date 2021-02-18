// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/BurntSushi/toml"
	"github.com/rande/goapp"
)

func LoadConfigurationFromFile(path string, c interface{}) error {
	data, err := goapp.LoadConfigurationFromFile(path)

	goapp.PanicOnError(err)

	return LoadConfiguration(data, c)
}

func LoadConfigurationFromString(conf string, c interface{}) error {
	data, err := goapp.LoadConfigurationFromString(conf)

	goapp.PanicOnError(err)

	return LoadConfiguration(data, c)
}

func LoadConfiguration(conf string, c interface{}) error {
	_, err := toml.Decode(conf, c)

	goapp.PanicOnError(err)

	return nil
}
