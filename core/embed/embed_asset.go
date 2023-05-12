// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

var contentTypes = map[string]string{
	"js":    "application/javascript",
	"css":   "text/css",
	"svg":   "image/svg+xml",
	"eot":   "application/vnd.ms-fontobject",
	"woff":  "application/x-font-woff",
	"woff2": "application/font-woff2",
	"ttf":   "application/x-font-ttf",
	"png":   "image/png",
	"jpg":   "image/jpg",
	"gif":   "image/gif",
}


func PageNotFound(res http.ResponseWriter) {
	res.WriteHeader(404)
		res.Write([]byte("<html><head><title>Page not found</title></head><body><h1>Page not found</h1></body></html>"))
}

func ConfigureEmbedMux(mux *web.Mux, embeds *Embeds, publicPath string, logger *log.Logger) {

	lenPath := len(publicPath)

	if logger != nil {
		logger.WithFields(log.Fields{
			"module":       "embed.mux",
			"public_path":  publicPath,
		}).Debug("Configure embed assets")
	}

	mux.Get(publicPath+"/*", func(c web.C, res http.ResponseWriter, req *http.Request) {
		var logger *log.Entry

		if l, ok := c.Env["logger"]; ok {
			logger = l.(*log.Entry).WithFields(log.Fields{
				"module": "embed.handler",
			})
		}

		path := req.RequestURI[lenPath:]

		if path[len(path)-1:] == "/" {
			path = path[0 : len(path)-1]
		}

		sections := strings.Split(path, ",")

		if (len(sections) < 2) {
			PageNotFound(res)
			return
		}

		module := sections[0]

		paths := []string{
			path,
			path + "/index.html",
		}

		for _, path := range paths {
			if logger != nil {
				logger.Debug("GET:", path)
			}

			asset, err := embeds.ReadFile(module, strings.Join(sections[1:], "/"))

			if err != nil {
				if logger != nil {
					logger.Debug("Err:", err)
				}

				continue
			}

			ext := path[(strings.LastIndex(path, ".") + 1):]

			if _, ok := contentTypes[ext]; ok {
				res.Header().Set("Content-Type", contentTypes[ext])

				logger.Debug("Content-Type:", contentTypes[ext])
			} else {
				res.Header().Set("Content-Type", "application/stream")
			}

			res.Write(asset)

			return
		}

		PageNotFound(res)
	})
}
