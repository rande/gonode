package extra

import (
	nc "github.com/rande/gonode/core"
	"github.com/BurntSushi/toml"
	"text/template"
	"os"
	"io/ioutil"
	"bytes"
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

	data, err := LoadConfigurationFromFile(path)

	nc.PanicOnError(err)

	_, err = toml.Decode(data, &config)

	nc.PanicOnError(err)

	return config
}


func LoadConfigurationFromFile(path string) (string, error) {

	data, err := ioutil.ReadFile(path)

	nc.PanicOnError(err)

	t := template.New("config")
	t.Funcs(map[string]interface {}{
		"env": os.Getenv,
	})
	_, err = t.Parse(string(data[:]))

	nc.PanicOnError(err)

	b := bytes.NewBuffer([]byte{})

	err = t.Execute(b, nil)

	nc.PanicOnError(err)

	return b.String(), nil
}
