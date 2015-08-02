// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package extra

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	rexMeta = regexp.MustCompile(`meta\.([a-zA-Z]*)`)
	rexData = regexp.MustCompile(`data\.([a-zA-Z]*)`)
)

func SendStatusMessage(res http.ResponseWriter, code int, message string) {
	res.WriteHeader(code)

	status := "KO"
	if code >= 200 && code < 300 {
		status = "OK"
	}

	data, _ := json.Marshal(map[string]string{
		"status":  status,
		"message": message,
	})

	res.Write(data)
}
