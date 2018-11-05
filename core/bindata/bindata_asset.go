// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	config "github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
)

var contentTypes = map[string]string{
	"js":    "application/javascript; charset=utf-8",
	"html":  "text/html; charset=utf-8",
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

func Send404(res http.ResponseWriter) {
	res.WriteHeader(404)
	res.Write([]byte("<html><head><title>Document not found</title></head><body><h1>Document not found</h1></body></html>"))
}

func ConfigureBinDataMux(mux *web.Mux, Asset func(name string) ([]byte, error), data *config.BinDataAsset, logger *log.Logger) {
	lenPath := len(data.Public)

	if logger != nil {
		logger.WithFields(log.Fields{
			"module":       "bindata.mux",
			"public_path":  data.Public,
			"private_path": data.Private,
			"fallback":     data.Fallback,
		}).Debug("Configure bindata assets")
	}

	mux.Get(data.Public+"/*", func(c web.C, res http.ResponseWriter, req *http.Request) {
		var logger *log.Entry

		if l, ok := c.Env["logger"]; ok {
			logger = l.(*log.Entry).WithFields(log.Fields{
				"module":       "bindata.handler",
				"public_path":  data.Public,
				"private_path": data.Private,
				"fallback":     data.Fallback,
			})
		}

		parsedURI, err := url.ParseRequestURI(req.RequestURI)

		if err != nil {
			if logger != nil {
				logger.WithFields(log.Fields{
					"request-uri": req.RequestURI,
				}).Warning("Unable to parse the url")
			}

			Send404(res)
			return
		}

		path := parsedURI.Path[lenPath:]

		if path[len(path)-1:] == "/" {
			path = path[0 : len(path)-1]
		}

		paths := []string{
			data.Private + path,
			data.Private + path + "/" + data.Index,
		}

		if len(data.Fallback) > 0 {
			paths = append(paths, data.Private+data.Fallback)
		}

		for _, path := range paths {
			asset, err := Asset(path)

			if err != nil {
				if logger != nil {
					logger.WithFields(log.Fields{
						"path": path,
					}).Debug("Unable to find asset")
				}

				continue
			}

			ext := path[(strings.LastIndex(path, ".") + 1):]

			contentType := "application/octet-stream"
			if _, ok := contentTypes[ext]; ok {
				contentType = contentTypes[ext]
			}

			if logger != nil {
				logger.WithFields(log.Fields{
					"content-type": contentType,
					"path":         path,
					"ext":          ext,
				}).Debug("Send asset contents")
			}

			res.Header().Set("Content-Type", contentType)
			res.Write(asset)

			return
		}

		Send404(res)
		return
	})
}
