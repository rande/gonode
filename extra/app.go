package extra

import (
	"database/sql"
	sq "github.com/lann/squirrel"
	pq "github.com/lib/pq"
	. "github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	nh "github.com/rande/gonode/handlers"
	"github.com/spf13/afero"
	"log"
	"net/http"
	"time"
)

func ConfigureApp(app *App) {
	app.Set("gonode.fs", func(app *App) interface{} {
		configuration := app.Get("gonode.configuration").(*Config)

		return nc.NewSecureFs(&afero.OsFs{}, configuration.Filesystem.Path)
	})

	app.Set("gonode.http_client", func(app *App) interface{} {
		return &http.Client{}
	})

	app.Set("gonode.handler_collection", func(app *App) interface{} {
		return nc.HandlerCollection{
			"default": &nh.DefaultHandler{},
			"media.image": &nh.ImageHandler{
				Fs: app.Get("gonode.fs").(*nc.SecureFs),
			},
			"media.youtube": &nh.YoutubeHandler{},
			"blog.post":     &nh.PostHandler{},
			"core.user":     &nh.UserHandler{},
		}
	})

	app.Set("gonode.manager", func(app *App) interface{} {
		return &nc.PgNodeManager{
			Logger:   app.Get("logger").(*log.Logger),
			Db:       app.Get("gonode.postgres.connection").(*sql.DB),
			ReadOnly: false,
			Handlers: app.Get("gonode.handler_collection").(nc.Handlers),
		}
	})

	app.Set("gonode.postgres.connection", func(app *App) interface{} {

		configuration := app.Get("gonode.configuration").(*Config)

		sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		db, err := sql.Open("postgres", configuration.Databases["master"].DSN)

		db.SetMaxIdleConns(8)
		db.SetMaxOpenConns(64)

		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		return db
	})

	app.Set("gonode.api", func(app *App) interface{} {
		return &nc.Api{
			Manager:    app.Get("gonode.manager").(*nc.PgNodeManager),
			Version:    "1.0.0",
			Serializer: app.Get("gonode.node.serializer").(*nc.Serializer),
		}
	})

	app.Set("gonode.node.serializer", func(app *App) interface{} {
		s := nc.NewSerializer()
		s.Handlers = app.Get("gonode.handler_collection").(nc.Handlers)

		return s
	})

	app.Set("gonode.postgres.subscriber", func(app *App) interface{} {
		configuration := app.Get("gonode.configuration").(*Config)

		return nc.NewSubscriber(
			configuration.Databases["master"].DSN,
			app.Get("logger").(*log.Logger),
		)
	})

	app.Set("gonode.listener.youtube", func(app *App) interface{} {
		client := app.Get("gonode.http_client").(*http.Client)

		return &nh.YoutubeListener{
			HttpClient: client,
		}
	})

	app.Set("gonode.listener.file_downloader", func(app *App) interface{} {
		client := app.Get("gonode.http_client").(*http.Client)
		fs := app.Get("gonode.fs").(*nc.SecureFs)

		return &nh.ImageDownloadListener{
			Fs:         fs,
			HttpClient: client,
		}
	})

	// need to find a way to trigger the handler registration
	sub := app.Get("gonode.postgres.subscriber").(*nc.Subscriber)

	sub.ListenMessage("media_youtube_update", func(app *App) nc.SubscriberHander {
		manager := app.Get("gonode.manager").(*nc.PgNodeManager)
		listener := app.Get("gonode.listener.youtube").(*nh.YoutubeListener)

		return func(notification *pq.Notification) (int, error) {
			return listener.Handle(notification, manager)
		}
	}(app))

	sub.ListenMessage("media_file_download", func(app *App) nc.SubscriberHander {
		manager := app.Get("gonode.manager").(*nc.PgNodeManager)
		listener := app.Get("gonode.listener.file_downloader").(*nh.ImageDownloadListener)

		return func(notification *pq.Notification) (int, error) {
			return listener.Handle(notification, manager)
		}
	}(app))

	sub.ListenMessage("core_sleep", func(app *App) nc.SubscriberHander {
		return func(notification *pq.Notification) (int, error) {

			logger := app.Get("logger").(*log.Logger)

			duration, _ := time.ParseDuration(notification.Extra)

			logger.Printf("[core_sleep] sleep ...")

			time.Sleep(duration)

			logger.Printf("[core_sleep] wake up ...")

			return nc.PubSubListenContinue, nil
		}
	}(app))
}
