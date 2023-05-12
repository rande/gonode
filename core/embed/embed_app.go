package embed

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("gonode.embeds", func(app *goapp.App) interface{} {
			return NewEmbeds()
		})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.pongo", func(app *goapp.App) interface{} {

			return pongo2.NewSet("gonode.embeds", &PongoTemplateLoader{
				Embeds: app.Get("gonode.embeds").(*Embeds),
				BasePath: "",
			})
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		if !app.Has("goji.mux") {
			return nil
		}

		// expose files using static/modules/[path]

		mux := app.Get("goji.mux").(*web.Mux)
		logger := app.Get("logger").(*log.Logger)
		asset := app.Get("gonode.embeds").(*Embeds)

		ConfigureEmbedMux(mux, asset, "/static", logger)

		return nil
	})
}
