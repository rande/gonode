// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/orandin/lumberjackrus"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/helper"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

var (
	MissingHookNameError = errors.New("missing hook name")
	NoHookHandlerError   = errors.New("No hook handler")
)

func GetValue(name string, conf map[string]interface{}, d ...interface{}) interface{} {
	if v, ok := conf[name]; !ok {
		if len(d) > 0 {
			return d[0]
		}

		panic(fmt.Sprintf("missing key: %s", name))
	} else {
		return v
	}
}

func GetHook(conf map[string]interface{}) (log.Hook, error) {
	if _, ok := conf["service"]; !ok {
		return nil, MissingHookNameError
	}

	var tags []string
	if _, ok := conf["tags"]; !ok {
		tags = nil
	} else {
		switch ts := conf["tags"].(type) {
		case []string:
			tags = ts
		case []interface{}:
			for _, tag := range ts {
				tags = append(tags, tag.(string))
			}
		default:
			panic("invalid type")
		}
	}

	switch conf["service"] {
	case "lumberjackrus":
		level, err := logrus.ParseLevel(GetValue("level", conf, "info").(string))

		if err != nil {
			return nil, err
		}

		hook, err := lumberjackrus.NewHook(
			&lumberjackrus.LogFile{
				Filename:   GetValue("filename", conf, fmt.Sprintf("%s/gonode.log", os.TempDir())).(string),
				MaxSize:    int(GetValue("max_size", conf, 100).(int64)),
				MaxBackups: int(GetValue("max_backups", conf, 0).(int)),
				MaxAge:     int(GetValue("max_age", conf, 31).(int64)),
				Compress:   GetValue("compress", conf, 100).(bool),
				LocalTime:  GetValue("local_time", conf, 100).(bool),
			},
			level,
			&logrus.TextFormatter{},
			&lumberjackrus.LogFileOpts{},
		)

		if err != nil {
			return nil, err
		}

		return hook, err
	}

	return nil, NoHookHandlerError
}

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("logger", func(app *goapp.App) interface{} {

			logger := log.New()
			logger.Out = os.Stdout
			logger.Level, _ = log.ParseLevel(strings.ToLower(conf.Logger.Level))

			d := &DispatchHook{
				make(map[log.Level][]log.Hook, 0),
			}

			for _, v := range conf.Logger.Hooks {
				hook, err := GetHook(v)

				helper.PanicOnError(err)

				l, err := log.ParseLevel(strings.ToLower(GetValue("level", v).(string)))

				helper.PanicOnError(err)

				d.Add(hook, l)
			}

			logger.Hooks.Add(d)

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
