// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/config"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"github.com/zenazn/goji/web/mutil"
	"net/http"
	"os"
	"time"
)

func GetMiddleware(logger *log.Logger) func(c *web.C, h http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(*c)

			fields := initFields(reqID, r)

			c.Env["logger"] = logger.WithFields(fields)

			lw := mutil.WrapWriter(w)

			t1 := time.Now()
			h.ServeHTTP(lw, r)

			if lw.Status() == 0 {
				lw.WriteHeader(http.StatusOK)
			}
			t2 := time.Now()

			fields = printEnd(fields, reqID, lw, t2.Sub(t1))

			logger.WithFields(fields).Info("Serve request")
		}

		return http.HandlerFunc(fn)
	}
}

func initFields(reqID string, r *http.Request) log.Fields {
	fields := log.Fields{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
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

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("logger", func(app *goapp.App) interface{} {

			logger := log.New()
			logger.Out = os.Stderr
			logger.Level = log.DebugLevel

			return logger
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		logger := app.Get("logger").(*log.Logger)

		mux := app.Get("goji.mux").(*web.Mux)

		mux.Use(GetMiddleware(logger))

		return nil
	})
}
