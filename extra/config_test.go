package extra

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/twinj/uuid"
	"os"
)

func Test_LoadConfigurationFromFile_WithEnv(t *testing.T) {
	os.Setenv("PG_USER", "foo")
	os.Setenv("PG_PASSWORD", "bar")

	defer func() {
		os.Unsetenv("PG_USER")
		os.Unsetenv("PG_PASSWORD")
	}()

	data, err := LoadConfigurationFromFile("../test/config_codeship.toml")

	expected := `name= "GoNode - Codeship"

[databases.master]
type    = "master"
dsn     = "postgres://foo:bar@localhost/test"
enabled = true
prefix  = "test"


[filesystem]
path = "/tmp/gnode"
`

	assert.Nil(t, err)
	assert.Equal(t, data, expected)
}

func Test_LoadConfigurationFromFile_WithoutEnv(t *testing.T) {
	data, err := LoadConfigurationFromFile("../test/config_codeship.toml")

	expected := `name= "GoNode - Codeship"

[databases.master]
type    = "master"
dsn     = "postgres://:@localhost/test"
enabled = true
prefix  = "test"


[filesystem]
path = "/tmp/gnode"
`
	assert.Nil(t, err)
	assert.Equal(t, data, expected)
}

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
	assert.Equal(t, config.Databases["master"].DSN, "postgres://foo:bar@localhost/test")
	assert.Equal(t, config.Databases["master"].Enabled, true)
	assert.Equal(t, config.Databases["master"].Prefix, "test")
	assert.Equal(t, config.Filesystem.Type, "") // not used for now
	assert.Equal(t, config.Filesystem.Path, "/tmp/gnode")
}
