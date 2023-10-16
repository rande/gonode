// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Execute_Success(t *testing.T) {

	user := &TestUser{
		Name:    "John Doe",
		Enabled: true,
		Hidden:  false,
	}

	form := CreateForm(user)
	form.Add("Name", "text")

	v := url.Values{
		"Name": []string{"Jane Doe"},
	}

	req, _ := http.NewRequest("POST", "/submit", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err := Process(form, req)

	assert.Nil(t, err)

}
