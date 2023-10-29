package embed

import (
	"fmt"
	"html/template"

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

			engine := pongo2.NewSet("gonode.embeds", &PongoTemplateLoader{
				Embeds:   app.Get("gonode.embeds").(*Embeds),
				BasePath: "",
			})

			engine.Options = &pongo2.Options{
				TrimBlocks:   true,
				LStripBlocks: true,
			}

			return engine
		})

		app.Set("gonode.template", func(app *goapp.App) interface{} {
			return &TemplateLoader{
				Embeds:   app.Get("gonode.embeds").(*Embeds),
				BasePath: "",
				Template: template.New("default"),
			}
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
		loader := app.Get("gonode.template").(*TemplateLoader)
		embeds := app.Get("gonode.embeds").(*Embeds)

		ConfigureTemplates(loader.Template, embeds)
		ConfigureEmbedMux(mux, asset, "/static", logger)

		return nil
	})
}

func ConfigureTemplates(tpl *template.Template, embeds *Embeds) error {
	entries := embeds.GetFilesByExt(".html")

	for _, entry := range entries {
		// reformat template name to respect the convention: module:template.html
		if len(entry.Path) < 10 {
			continue
		}

		name := fmt.Sprintf("%s:%s", entry.Module, entry.Path[10:])

		data, err := embeds.ReadFile(entry.Module, entry.Path)

		if err != nil {
			fmt.Printf("Unable to read file: %s\n", err)
			return err
		}

		_, err = tpl.New(name).Parse(string(data))

		if err != nil {
			fmt.Printf("Error parsing the template: %s, %s\n", name, err)
			return err
		}
	}

	return nil
}
