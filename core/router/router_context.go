// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"strings"
)

func RequestContextMiddleware(c *web.C, h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context, _ := BuildRequestContext(r)

		c.Env["request_context"] = context

		h.ServeHTTP(w, r)
	})
}

type RequestContext struct {
	Host     string
	Port     int
	Protocol string
	Prefix   string
}

func BuildRequestContext(req *http.Request) (*RequestContext, error) {
	c := &RequestContext{}

	// resolve host
	c.Host = req.Host
	c.Protocol = "http"
	c.Port = 80

	if len(req.Header.Get("X-Forwarded-Proto")) > 0 {
		c.Protocol = req.Header.Get("X-Forwarded-Proto")

		if c.Protocol == "https" {
			c.Port = 443
		}
	}

	if len(req.Header.Get("X-Forwarded-Host")) > 0 {
		c.Host = req.Header.Get("X-Forwarded-Host")
	}

	sep := strings.LastIndex(c.Host, ":")
	if sep > 0 {
		port, err := strconv.Atoi(c.Host[sep+1:])

		if err != nil {
			return nil, err
		}

		c.Host = c.Host[:sep]
		c.Port = port
	}

	c.Prefix = c.Protocol + "://" + c.Host

	if (c.Port != 80 && c.Protocol == "http") || (c.Port != 443 && c.Protocol == "https") {
		c.Prefix = fmt.Sprintf("%s:%d", c.Prefix, c.Port)
	}

	return c, nil
}
