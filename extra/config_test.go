// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package extra

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/twinj/uuid"
	"os"
)

func Test_GetConfiguration(t *testing.T) {
	os.Setenv("PG_USER", "foo")
	os.Setenv("PG_PASSWORD", "bar")

	defer func() {
		os.Unsetenv("PG_USER")
		os.Unsetenv("PG_PASSWORD")
	}()

	config := GetConfiguration("../test/config_codeship.toml")

	assert.Equal(t, config.Name, "GoNode - Codeship")
	assert.Equal(t, config.Databases["master"].Type, "master")
	assert.Equal(t, config.Databases["master"].DSN, "postgres://foo:bar@localhost:5434/test")
	assert.Equal(t, config.Databases["master"].Enabled, true)
	assert.Equal(t, config.Databases["master"].Prefix, "test")
	assert.Equal(t, config.Filesystem.Type, "") // not used for now
	assert.Equal(t, config.Filesystem.Path, "/tmp/gnode")
}
