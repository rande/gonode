// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/BurntSushi/toml"
	"github.com/rande/goapp"
)

type ServerDatabase struct {
	Name    string `toml:"name"`
	DSN     string `toml:"dsn"`
	Type    string `toml:"type"`
	Prefix  string `toml:"prefix"`
	Enabled bool   `toml:"enabled"`
}

type ServerFilesystem struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

type ServerHandler struct {
	Type    string `toml:"type"`
	Enabled bool   `toml:"enabled"`
}

type ServerConfig struct {
	Name       string                     `toml:"name"`
	Databases  map[string]*ServerDatabase `toml:"databases"`
	Filesystem ServerFilesystem           `toml:"filesystem"`
	Test       bool                       `toml:"test"`
	Bind       string                     `toml:"bind"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Databases: make(map[string]*ServerDatabase),
		Bind:      ":2408",
		Test:      false,
	}
}

func LoadConfiguration(path string, c interface{}) error {
	data, err := goapp.LoadConfigurationFromFile(path)

	goapp.PanicOnError(err)

	_, err = toml.Decode(data, c)

	goapp.PanicOnError(err)

	return nil
}
