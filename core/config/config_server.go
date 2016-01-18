// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"os"
)

type ServerSearch struct {
	MaxResult uint64 `toml:"max_result"`
}

type ServerGuard struct {
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

type ServerBinDataAsset struct {
	Index   string `toml:"index"`
	Public  string `toml:"public"`
	Private string `toml:"private"`
}

type ServerBinData struct {
	BasePath  string                         `toml:"base_path"`
	Assets    map[string]*ServerBinDataAsset `toml:"assets"`
	Templates []string                       `toml:"templates"`
}

type ServerSecurity struct {
	Cors struct {
		AllowedOrigins     []string `toml:"allowed_origins"`
		AllowedMethods     []string `toml:"allowed_methods"`
		AllowedHeaders     []string `toml:"allowed_headers"`
		ExposedHeaders     []string `toml:"exposes_headers"`
		AllowCredentials   bool     `toml:"allow_credentials"`
		MaxAge             int      `toml:"max_age"`
		OptionsPassthrough bool     `toml:"options_passthrough"`
	} `toml:"cors"`
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
	Guard      *ServerGuard               `toml:"guard"`
	Security   *ServerSecurity            `toml:"security"`
	Search     *ServerSearch              `toml:"search"`
	BinData    *ServerBinData             `toml:"bindata"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Databases: make(map[string]*ServerDatabase),
		Bind:      ":2408",
		Test:      false,
		Search: &ServerSearch{
			MaxResult: 128,
		},
		BinData: &ServerBinData{
			BasePath: os.Getenv("GOPATH") + "/src",
			Assets:   make(map[string]*ServerBinDataAsset, 0),
			//			Templates: make([]string, 0),
		},
	}
}
