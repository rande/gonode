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

		loader.Templates = GetTemplates(embeds)

		ConfigureEmbedMux(mux, asset, "/static", logger)

		return nil
	})
}

// This function is called only once at boot time to configure the different template
func GetTemplates(embeds *Embeds) map[string]*template.Template {
	entries := embeds.GetFilesByExt(".html")
	// in the entries we need to find the page, each page will have its own set of templates (layout, blocks, etc ...)

	templates := map[string]*template.Template{}

	// create root template without parsing them
	pagesPath := "templates/pages/"
	for _, entry := range entries {
		if len(entry.Path) < len(pagesPath) || entry.Path[0:len(pagesPath)] != pagesPath {
			continue
		}

		name := fmt.Sprintf("%s:%s", entry.Module, entry.Path[10:len(entry.Path)-5])
		templates[name] = template.New(name)
	}

	layoutsPath := "templates/layouts/"
	blocksPath := "templates/blocks/"

	// load all the layout first, default templates will be defined
	for name, tpl := range templates {

		fmt.Printf("Iterating over layout: %s\n", name)
		for _, entry := range entries {
			if len(entry.Path) < len(layoutsPath) || entry.Path[0:len(layoutsPath)] != layoutsPath {
				continue
			}

			name := fmt.Sprintf("%s:%s", entry.Module, entry.Path[10:len(entry.Path)-5])

			fmt.Printf("Loading layout: %s\n", name)

			if data, err := embeds.ReadFile(entry.Module, entry.Path); err != nil {
				fmt.Printf("Unable to read file: %s\n", err)
				panic(err)
			} else if _, err = tpl.New(name).Parse(string(data)); err != nil {
				fmt.Printf("Error parsing the template: %s, %s\n", name, err)
				panic(err)
			}
		}

		// load all the blocks first, so this will let an option to overwrite them if needed in
		// the page

		fmt.Printf("Iterating over block: %s\n", name)
		for _, entry := range entries {
			if len(entry.Path) < len(blocksPath) || entry.Path[0:len(blocksPath)] != blocksPath {
				continue
			}

			name := fmt.Sprintf("%s:%s", entry.Module, entry.Path[10:len(entry.Path)-5])

			fmt.Printf("Loading blocks: %s\n", name)

			if data, err := embeds.ReadFile(entry.Module, entry.Path); err != nil {
				fmt.Printf("Unable to read file: %s\n", err)
				panic(err)
			} else if _, err = tpl.New(name).Parse(string(data)); err != nil {
				fmt.Printf("Error parsing the template: %s, %s\n", name, err)
				panic(err)
			}
		}
	}

	// we need to load the main template last in order to ensure the defined template will be
	// the last registered in the template stack.
	for _, entry := range entries {
		if len(entry.Path) < len(pagesPath) || entry.Path[0:len(pagesPath)] != pagesPath {
			continue
		}

		name := fmt.Sprintf("%s:%s", entry.Module, entry.Path[10:len(entry.Path)-5])
		if data, err := embeds.ReadFile(entry.Module, entry.Path); err != nil {
			fmt.Printf("Unable to read file: %s\n", err)
			panic(err)
		} else {
			templates[name].Parse(string(data))
		}
	}

	return templates
}
