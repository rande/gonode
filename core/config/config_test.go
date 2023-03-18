// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func Test_Server_LoadConfigurationFromFile(t *testing.T) {
	os.Setenv("PG_USER", "foo")
	os.Setenv("PG_PASSWORD", "bar")

	defer func() {
		os.Unsetenv("PG_USER")
		os.Unsetenv("PG_PASSWORD")
	}()

	config := &Config{
		Databases: make(map[string]*Database),
	}

	err := LoadConfigurationFromString(`
name= "GoNode - Codeship"
bind= ":2508"

[databases.master]
type    = "master"
dsn     = "postgres://{{ env "PG_USER" }}:{{ env "PG_PASSWORD" }}@localhost:5434/test"
enabled = true
prefix  = "test"


[filesystem]
path = "/tmp/gnode"

[guard]
key = "ZeSecretKey0oo"

    [guard.jwt]
        [guard.jwt.login]
        path = "/login"

        [guard.jwt.token]
        path = "^\\/nodes\\/(.*)$"

[security]
    voters = ["security.voter.role"]

    [security.cors]
    allowed_origins = ["*"]
    allowed_methods = ["GET", "PUT", "POST"]
    allowed_headers = ["Origin", "Accept", "Content-Type", "Authorization"]

    [[security.access]]
    path  = "^\\/admin"
    roles = ["ROLE_ADMIN"]

[search]
    max_result = 256

[bindata]
    base_path = "/var/go"
    templates = [
        "/path/post/templates",
        "/path/media/templates",
        #"/commented"
    ]

    [bindata.assets]
        [bindata.assets.explorer]
        public = "/explorer"
        private = "github.com/rande/gonode/explorer/dist"

[media]
    [media.image]
    allowed_widths = [100, 200]
    max_width = 300

[logger]

    level = "debug"

    [logger.fields]
    app = "gonode"

`, config)

	assert.NoError(t, err)

	// test general configuration
	assert.Equal(t, config.Name, "GoNode - Codeship")
	assert.Equal(t, config.Databases["master"].Type, "master")
	assert.Equal(t, config.Databases["master"].DSN, "postgres://foo:bar@localhost:5434/test")
	assert.Equal(t, config.Databases["master"].Enabled, true)
	assert.Equal(t, config.Databases["master"].Prefix, "test")
	assert.Equal(t, config.Filesystem.Type, "") // not used for now
	assert.Equal(t, config.Filesystem.Path, "/tmp/gnode")

	// test guard
	assert.Equal(t, config.Guard.Jwt.Login.Path, "/login")
	assert.Equal(t, config.Guard.Jwt.Token.Path, `^\/nodes\/(.*)$`)

	// test security: cors
	assert.False(t, config.Security.Cors.AllowCredentials)
	assert.Equal(t, config.Security.Cors.AllowedHeaders, []string{"Origin", "Accept", "Content-Type", "Authorization"})
	assert.Equal(t, config.Security.Cors.AllowedMethods, []string{"GET", "PUT", "POST"})

	// test security: access
	assert.Equal(t, 1, len(config.Security.Access))
	assert.Equal(t, []string{"ROLE_ADMIN"}, config.Security.Access[0].Roles)
	assert.Equal(t, "^\\/admin", config.Security.Access[0].Path)
	assert.Equal(t, []string{"security.voter.role"}, config.Security.Voters)

	// test search
	assert.Equal(t, uint64(256), config.Search.MaxResult)

	// test bindata
	assert.Equal(t, "/explorer", config.BinData.Assets["explorer"].Public)
	assert.Equal(t, "github.com/rande/gonode/explorer/dist", config.BinData.Assets["explorer"].Private)
	assert.Equal(t, "/var/go", config.BinData.BasePath)
	assert.Equal(t, []string{"/path/post/templates", "/path/media/templates"}, config.BinData.Templates)

	// test media
	assert.Equal(t, uint(300), config.Media.Image.MaxWidth)
	assert.Equal(t, []uint{100, 200}, config.Media.Image.AllowedWidths)

	// test logger
	assert.Equal(t, map[string]string{"app": "gonode"}, config.Logger.Fields)

	// debug
	config.Guard.Jwt.Login.Path = `^\/nodes\/(.*)$`

	w := bytes.NewBufferString("")
	e := toml.NewEncoder(w)

	e.Encode(config)
}
