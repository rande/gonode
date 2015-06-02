package extra

import (
	"github.com/BurntSushi/toml"
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

	_, err := toml.DecodeFile(path, &config)

	if err != nil {
		panic(err)
	}

	return config
}
