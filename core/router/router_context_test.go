// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RequestContext_Without_X(t *testing.T) {
	req, _ := http.NewRequest("GET", "/Hello", bytes.NewBuffer([]byte("")))
	req.Host = "gonode.com:80"

	c, err := BuildRequestContext(req)

	assert.NoError(t, err)
	assert.Equal(t, "gonode.com", c.Host)
	assert.Equal(t, 80, c.Port)
	assert.Equal(t, "http://gonode.com", c.Prefix)
}

func Test_RequestContext_With_X(t *testing.T) {
	req, _ := http.NewRequest("GET", "/Hello", bytes.NewBuffer([]byte("")))
	req.Host = "localhost:80"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "gonode.com:443")

	c, err := BuildRequestContext(req)

	assert.NoError(t, err)
	assert.Equal(t, "gonode.com", c.Host)
	assert.Equal(t, 443, c.Port)
	assert.Equal(t, "https://gonode.com", c.Prefix)
}

func Test_RequestContext_With_X_NonDefaultPort(t *testing.T) {
	req, _ := http.NewRequest("GET", "/Hello", bytes.NewBuffer([]byte("")))
	req.Host = "localhost:80"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "gonode.com:445")

	c, err := BuildRequestContext(req)

	assert.NoError(t, err)
	assert.Equal(t, "gonode.com", c.Host)
	assert.Equal(t, 445, c.Port)
	assert.Equal(t, "https://gonode.com:445", c.Prefix)
}
