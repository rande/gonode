// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type AccessRule struct {
	Path  *regexp.Regexp
	Roles Attributes
}

type AccessChecker struct {
	Rules         []*AccessRule
	DecisionVoter DecisionVoter
}

func (c *AccessChecker) Check(t SecurityToken, req *http.Request) bool {

	roles := c.getRoles(req)

	return c.DecisionVoter.Decide(t, roles, req)
}

func (c *AccessChecker) getRoles(req *http.Request) Attributes {
	attrs := make(Attributes, 0)

	for _, r := range c.Rules {
		if r.Path.Match([]byte(req.URL.Path)) {
			return r.Roles
		}
	}

	return attrs
}

func RenderForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("<html><head><title>Forbidden</title></head><body><h1>Forbidden</h1><p>Sorry, you don't have enough credentials to view the requested resource.</p></body></html>"))
}

func AccessCheckerMiddleware(ac *AccessChecker) func(c *web.C, h http.Handler) http.Handler {

	return func(c *web.C, h http.Handler) http.Handler {
		var logger *log.Entry
		var token SecurityToken

		fn := func(w http.ResponseWriter, r *http.Request) {

			if _, ok := c.Env["logger"]; ok {
				logger = c.Env["logger"].(*log.Entry)
			}

			if _, ok := c.Env["guard_token"]; ok {
				token = c.Env["guard_token"].(SecurityToken)
			} else {
				// no token available, this is an error
				if logger != nil {
					logger.WithFields(log.Fields{
						"module": "core.security.access.middleware",
					}).Warn("Access Forbidden: no security token found")
				}

				RenderForbidden(w)

				return
			}

			if !ac.Check(token, r) {
				if logger != nil {
					logger.WithFields(log.Fields{
						"module":   "core.security.access.middleware",
						"roles":    token.GetRoles(),
						"username": token.GetUsername(),
					}).Warn("Access Forbidden: no enough credentials found")
				}

				RenderForbidden(w)

				return
			}

			if logger != nil {
				logger.WithFields(log.Fields{
					"module":   "core.security.access.middleware",
					"roles":    token.GetRoles(),
					"username": token.GetUsername(),
				}).Debug("Access granted")
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
