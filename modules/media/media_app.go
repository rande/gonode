// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	log "github.com/Sirupsen/logrus"
	pq "github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"net/http"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.listener.youtube", func(app *goapp.App) interface{} {
			return &YoutubeListener{
				HttpClient: app.Get("gonode.http_client").(*http.Client),
			}
		})

		app.Set("gonode.listener.file_downloader", func(app *goapp.App) interface{} {
			return &ImageDownloadListener{
				Vault:      app.Get("gonode.vault.fs").(*vault.Vault),
				HttpClient: app.Get("gonode.http_client").(*http.Client),
				Logger:     app.Get("logger").(*log.Logger),
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {

		c := app.Get("gonode.handler_collection").(base.HandlerCollection)
		c.Add("media.image", &ImageHandler{
			Vault:  app.Get("gonode.vault.fs").(*vault.Vault),
			Logger: app.Get("logger").(*log.Logger),
		})
		c.Add("media.youtube", &YoutubeHandler{})

		cv := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
		cv.Add("media.image", &MediaViewHandler{
			Vault:         app.Get("gonode.vault.fs").(*vault.Vault),
			MaxWidth:      conf.Media.Image.MaxWidth,
			AllowedWidths: conf.Media.Image.AllowedWidths,
		})

		// need to find a way to trigger the handler registration
		sub := app.Get("gonode.postgres.subscriber").(*base.Subscriber)

		sub.ListenMessage("media_youtube_update", func(app *goapp.App) base.SubscriberHander {
			manager := app.Get("gonode.manager").(*base.PgNodeManager)
			listener := app.Get("gonode.listener.youtube").(*YoutubeListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("media_file_download", func(app *goapp.App) base.SubscriberHander {
			manager := app.Get("gonode.manager").(*base.PgNodeManager)
			listener := app.Get("gonode.listener.file_downloader").(*ImageDownloadListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		return nil
	})

}
