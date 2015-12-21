// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Server_LoadConfigurationFromFile(t *testing.T) {
	os.Setenv("PG_USER", "foo")
	os.Setenv("PG_PASSWORD", "bar")

	defer func() {
		os.Unsetenv("PG_USER")
		os.Unsetenv("PG_PASSWORD")
	}()

	config := &ServerConfig{
		Databases: make(map[string]*ServerDatabase),
	}

	LoadConfigurationFromString(`
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
    [security.cors]
    allowed_origins = ["*"]
    allowed_methods = ["GET", "PUT", "POST"]
    allowed_headers = ["Origin", "Accept", "Content-Type", "Authorization"]

`, config)

	assert.Equal(t, config.Name, "GoNode - Codeship")
	assert.Equal(t, config.Databases["master"].Type, "master")
	assert.Equal(t, config.Databases["master"].DSN, "postgres://foo:bar@localhost:5434/test")
	assert.Equal(t, config.Databases["master"].Enabled, true)
	assert.Equal(t, config.Databases["master"].Prefix, "test")
	assert.Equal(t, config.Filesystem.Type, "") // not used for now
	assert.Equal(t, config.Filesystem.Path, "/tmp/gnode")

	assert.Equal(t, config.Guard.Jwt.Login.Path, "/login")
	assert.Equal(t, config.Guard.Jwt.Token.Path, `^\/nodes\/(.*)$`)

	assert.False(t, config.Security.Cors.AllowCredentials)
	assert.Equal(t, config.Security.Cors.AllowedHeaders, []string{"Origin", "Accept", "Content-Type", "Authorization"})
	assert.Equal(t, config.Security.Cors.AllowedMethods, []string{"GET", "PUT", "POST"})

	config.Guard.Jwt.Login.Path = `^\/nodes\/(.*)$`

	w := bytes.NewBufferString("")
	e := toml.NewEncoder(w)

	e.Encode(config)
}
