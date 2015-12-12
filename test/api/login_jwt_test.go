// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/plugins/user"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Create_Username(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		configuration := app.Get("gonode.configuration").(*config.ServerConfig)

		// WITH
		// create a valid user into the database ...
		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		u := app.Get("gonode.handler_collection").(core.HandlerCollection).NewNode("core.user")
		data := u.Data.(*user.User)
		data.Email = "test@example.org"
		data.Enabled = true
		data.FirstName = "Thomas"
		data.LastName = "Rxxxx"
		data.NewPassword = "ZePassword"
		data.Username = "rande"

		meta := u.Meta.(*user.UserMeta)
		meta.PasswordCost = 1 // save test time

		manager.Save(u, false)

		res, _ := test.RunRequest("POST", fmt.Sprintf("%s/login", ts.URL), url.Values{
			"username": {data.Username},
			"password": {"ZePassword"},
		})

		assert.Equal(t, 200, res.StatusCode)

		b := bytes.NewBuffer([]byte(""))
		io.Copy(b, res.Body)

		token, err := jwt.Parse(b.String(), func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(configuration.Auth.Key), nil
		})

		assert.NotNil(t, configuration.Auth.Key)
		assert.True(t, len(configuration.Auth.Key) > 0)
		assert.Nil(t, err)
		assert.True(t, token.Valid)

		fmt.Printf("%v", token)
	})
}
