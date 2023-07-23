// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"github.com/zenazn/goji/web/mutil"
)

func GetLoggerFromContext(c web.C) *log.Entry {
	if logger, ok := c.Env["logger"]; ok {
		return logger.(*log.Entry)
	}

	return nil
}

func GetMiddleware(logger *log.Logger) func(c *web.C, h http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(*c)

			fields := initFields(reqID, r)

			logger.WithFields(fields).Debug("Start serving a new request")

			c.Env["logger"] = logger.WithFields(fields)

			lw := mutil.WrapWriter(w)

			t1 := time.Now()
			h.ServeHTTP(lw, r)

			if lw.Status() == 0 {
				lw.WriteHeader(http.StatusOK)
			}
			t2 := time.Now()

			fields = printEnd(fields, reqID, lw, t2.Sub(t1))

			logger.WithFields(fields).Debug("Serve request")
		}

		return http.HandlerFunc(fn)
	}
}

func initFields(reqID string, r *http.Request) log.Fields {
	fields := log.Fields{
		"request_method":      r.Method,
		"request_url":         r.URL.String(),
		"request_remote_addr": r.RemoteAddr,
		"request_host":        r.Host,
	}

	if reqID != "" {
		fields["request_id"] = reqID
	}

	return fields
}

func printEnd(fields log.Fields, reqID string, w mutil.WriterProxy, dt time.Duration) log.Fields {
	if reqID != "" {
		fields["request_id"] = reqID
	}

	fields["status"] = w.Status()
	fields["content_type"] = w.Header().Get("Content-Type")
	fields["time"] = dt.Nanoseconds() / (int64(time.Millisecond) / int64(time.Nanosecond))

	return fields
}
