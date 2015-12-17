// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

type ServerAuth struct {
	Key string `toml:"key"`
	Jwt struct {
		Validity int64 `toml:"validity"`
		Login    struct {
			Path string `toml:"path"`
		} `toml:"login"`
		Token struct {
			Path string `toml:"path"`
		} `toml:"token"`
	} `toml:"jwt"`
}

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
	Auth       ServerAuth                 `toml:"auth"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Databases: make(map[string]*ServerDatabase),
		Bind:      ":2408",
		Test:      false,
	}
}
