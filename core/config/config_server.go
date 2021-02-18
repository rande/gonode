// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"os"
)

type Search struct {
	MaxResult uint64 `toml:"max_result"`
}

type Guard struct {
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
	Anonymous struct {
		Roles []string `toml:"roles"`
	}
}

type BinDataAsset struct {
	Index   string `toml:"index"`
	Public  string `toml:"public"`
	Private string `toml:"private"`
}

type BinData struct {
	BasePath  string                   `toml:"base_path"`
	Assets    map[string]*BinDataAsset `toml:"assets"`
	Templates []string                 `toml:"templates"`
}

type Security struct {
	Cors struct {
		AllowedOrigins     []string `toml:"allowed_origins"`
		AllowedMethods     []string `toml:"allowed_methods"`
		AllowedHeaders     []string `toml:"allowed_headers"`
		ExposedHeaders     []string `toml:"exposed_headers"`
		AllowCredentials   bool     `toml:"allow_credentials"`
		MaxAge             int      `toml:"max_age"`
		OptionsPassthrough bool     `toml:"options_passthrough"`
	} `toml:"cors"`
	Access []*struct {
		Path  string   `toml:"path"`
		Roles []string `toml:"roles"`
	} `toml:"access"`
	Voters []string `toml:"voters"`
}

type MediaImage struct {
	AllowedWidths []uint `toml:"allowed_widths"`
	MaxWidth      uint   `toml:"max_width"`
}

type Media struct {
	Image *MediaImage `toml:"image"`
}

type Database struct {
	Name    string `toml:"name"`
	DSN     string `toml:"dsn"`
	Type    string `toml:"type"`
	Prefix  string `toml:"prefix"`
	Enabled bool   `toml:"enabled"`
}

type Filesystem struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

type Handler struct {
	Type    string `toml:"type"`
	Enabled bool   `toml:"enabled"`
}

type Api struct {
	Prefix string `toml:"prefix"`
}

type Logger struct {
	Level  string                            `toml:"level"`
	Fields map[string]string                 `toml:"fields"`
	Hooks  map[string]map[string]interface{} `toml:"hooks"`
}

type Config struct {
	Name       string               `toml:"name"`
	Databases  map[string]*Database `toml:"databases"`
	Filesystem Filesystem           `toml:"filesystem"`
	Test       bool                 `toml:"test"`
	Bind       string               `toml:"bind"`
	Guard      *Guard               `toml:"guard"`
	Security   *Security            `toml:"security"`
	Search     *Search              `toml:"search"`
	BinData    *BinData             `toml:"bindata"`
	Media      *Media               `toml:"media"`
	Logger     *Logger              `toml:"logger"`
	Api        *Api                 `toml:"api"`
}

func NewConfig() *Config {
	return &Config{
		Databases: make(map[string]*Database),
		Bind:      ":2408",
		Test:      false,
		Search: &Search{
			MaxResult: 128,
		},
		BinData: &BinData{
			BasePath: os.Getenv("GOPATH") + "/src",
			Assets:   make(map[string]*BinDataAsset, 0),
		},
		Media: &Media{
			Image: &MediaImage{
				MaxWidth: uint(1024),
			},
		},
		Logger: &Logger{
			Level: "warn",
		},
		Api: &Api{
			Prefix: "/api",
		},
	}
}
