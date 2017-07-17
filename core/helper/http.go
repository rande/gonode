// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"encoding/json"
	"net/http"
)

func SendWithHttpCode(res http.ResponseWriter, code int, message string) {
	res.Header().Set("Content-Type", "application/json")

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

func SendWithStatus(status string, message string, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")

	if status == "KO" {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}

	data, _ := json.Marshal(map[string]string{
		"status":  status,
		"message": message,
	})

	res.Write(data)
}
