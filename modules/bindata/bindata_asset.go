// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rande/gonode/assets"
	"github.com/zenazn/goji/web"
	"net/http"
	"strings"
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

func ConfigureBinDataMux(mux *web.Mux, publicPath, privatePath, index string, logger *log.Logger) {

	lenPath := len(publicPath)

	mux.Get(publicPath+"/*", func(c web.C, res http.ResponseWriter, req *http.Request) {
		var logger *log.Entry

		if l, ok := c.Env["logger"]; ok {
			logger = l.(*log.Entry).WithFields(log.Fields{
				"module": "bindata.handler",
			})
		}

		path := req.RequestURI[lenPath:]

		if path[len(path)-1:] == "/" {
			path = path[0 : len(path)-1]
		}

		paths := []string{
			privatePath + path,
			privatePath + path + "/" + index,
		}

		for _, path := range paths {
			if logger != nil {
				logger.Debug("GET:", path)
			}

			asset, err := assets.Asset(path)

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
			}

			res.Write(asset)

			return
		}

		res.WriteHeader(404)
		res.Write([]byte("<html><head><title>Page not found</title></head><body><h1>Page not found</h1></body></html>"))

		return
	})
}
