// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package extra

import (
	"github.com/BurntSushi/toml"
	"github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
)

type Database struct {
	Name    string
	DSN     string
	Type    string
	Prefix  string
	Enabled bool
}

type Filesystem struct {
	Type string
	Path string
}

type Handler struct {
	Type    string
	Enabled bool
}

type Config struct {
	Name       string
	Databases  map[string]*Database
	Filesystem Filesystem
}

func GetConfiguration(path string) *Config {
	config := &Config{
		Databases: make(map[string]*Database),
	}

	data, err := goapp.LoadConfigurationFromFile(path)

	nc.PanicOnError(err)

	_, err = toml.Decode(data, &config)

	nc.PanicOnError(err)

	return config
}
