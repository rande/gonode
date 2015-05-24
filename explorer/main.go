package main

import (
  nc "github.com/rande/gonode/core"
  . "github.com/rande/goapp"
  nh "github.com/rande/gonode/handlers"
  "github.com/rande/gonode/extra"
  "net/http"
  pq "github.com/lib/pq"
  "database/sql"
  "log"
  sq "github.com/lann/squirrel"
  "github.com/spf13/afero"
  "strconv"
  "os"
  "github.com/hypebeast/gojistatic"
  "github.com/zenazn/goji/web"
  "github.com/zenazn/goji/web/middleware"
  "github.com/zenazn/goji/bind"
  "github.com/zenazn/goji/graceful"
  //  "time"
  "flag"
  "encoding/json"
)

func init() {
  bind.WithFlag()
  if fl := log.Flags(); fl&log.Ltime != 0 {
    log.SetFlags(fl | log.Lmicroseconds)
  }
  //  graceful.DoubleKickWindow(2 * time.Second)
}

func Serve(mux *web.Mux) {
  if !flag.Parsed() {
    flag.Parse()
  }

  mux.Compile()
  // Install our handler at the root of the standard net/http default mux.
  // This allows packages like expvar to continue working as expected.
  http.Handle("/", mux)

  listener := bind.Default()
  log.Println("Starting Goji on", listener.Addr())

  graceful.HandleSignals()
  bind.Ready()
  graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
  graceful.PostHook(func() { log.Printf("Goji stopped") })

  err := graceful.Serve(listener, http.DefaultServeMux)

  if err != nil {
    log.Fatal(err)
  }

  graceful.Wait()
}

func main() {
  //	LoadFixtures(manager, 256)
  //	Check(manager, manager.Logger)

  app := NewApp()

  // TODO: move this code to a dedicated init method in the extra folder
  //       to share common code
  // configure main services
  app.Set("logger", func(app *App) interface{} {
      return log.New(os.Stdout, "", log.Lshortfile)
    })

  app.Set("gonode.fs", func(app *App) interface{} {
      return nc.NewSecureFs(&afero.OsFs{}, "/tmp/gnode")
    })

  app.Set("gonode.http_client", func(app *App) interface{} {
      return &http.Client{}
    })

  app.Set("gonode.manager", func(app *App) interface{} {
      fs := app.Get("gonode.fs").(*nc.SecureFs)

      return &nc.PgNodeManager{
      Logger: app.Get("logger").(*log.Logger),
      Db: app.Get("gonode.postgres.connection").(*sql.DB),
      ReadOnly: false,
      Handlers: map[string]nc.Handler{
      "default":       &nh.DefaultHandler{},
      "media.image":   &nh.ImageHandler{
      Fs: fs,
  },
      "media.youtube": &nh.YoutubeHandler{},
      "blog.post":     &nh.PostHandler{},
      "core.user":     &nh.UserHandler{},
  },
  }
    })

  app.Set("gonode.postgres.connection", func(app *App) interface{} {
      sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
      db, err := sql.Open("postgres", "postgres://safre:safre@localhost/safre")
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
      Manager: app.Get("gonode.manager").(*nc.PgNodeManager),
      Version: "1.0.0",
  }
    })

  app.Set("goji.mux", func(app *App) interface{} {
      mux := web.New()

      mux.Use(middleware.RequestID)
      mux.Use(middleware.Logger)
      mux.Use(middleware.Recoverer)
      mux.Use(middleware.AutomaticOptions)

      mux.Use(gojistatic.Static("dist", gojistatic.StaticOptions{SkipLogging: true, Prefix: "dist"}))

      return mux
    })

  app.Set("gonode.postgres.subscriber", func(app *App) interface{} {
      return nc.NewSubscriber(
        "postgres://safre:safre@localhost/safre",
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
      Fs: fs,
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

  // load the current application
  extra.ConfigureGoji(app)

  ConfigureGoji(app)

  // start http server
  Serve(app.Get("goji.mux").(*web.Mux))
}


func Check(manager *nc.PgNodeManager, logger *log.Logger) {
  query := manager.SelectBuilder()
  nodes := manager.FindBy(query, 0, 10)

  logger.Printf("[PgNode] Found: %d nodes\n", nodes.Len())

  for e := nodes.Front(); e != nil; e = e.Next() {
    node := e.Value.(*nc.Node)
    logger.Printf("[PgNode] %s => %s", node.Uuid, node.Type)
  }

  node := manager.FindOneBy(manager.SelectBuilder().Where("type = ?", "media.image"))

  logger.Printf("[PgNode] Get first node: %s => %s (id: %d)", node.Uuid, node.Type, node.Id())

  node = manager.Find(node.Uuid)

  logger.Printf("[PgNode] Reload node, %s => %s (id: %d)", node.Uuid, node.Type, node.Id())

  node.Name = "This is my name"
  node.Slug = "this-is-my-name"

  manager.Save(node)

  manager.RemoveOne(node)

  manager.Remove(manager.SelectBuilder().Where("type = ?", "media.image"))

  logger.Printf("%s\n", "End code")
}

func Send(status string, message string, res http.ResponseWriter) {
  res.Header().Set("Content-Type", "application/json")

  if status == "KO" {
    res.WriteHeader(http.StatusInternalServerError)
  } else {
    res.WriteHeader(http.StatusOK)
  }

  data, _ := json.Marshal(map[string]string{
      "status": status,
      "message": message,
  })

  res.Write(data)
}

func ConfigureGoji(app *App) {

  mux := app.Get("goji.mux").(*web.Mux)
  manager := app.Get("gonode.manager").(*nc.PgNodeManager)
  prefix := ""

  mux.Put(prefix+"/data/purge", func(res http.ResponseWriter, req *http.Request) {

      tx, _ := manager.Db.Begin()
      manager.Db.Exec(`DELETE FROM "public"."nodes"`)
      manager.Db.Exec(`DELETE FROM "public"."nodes_audit"`)

      err := tx.Commit()

      if err != nil {
        Send("KO", err.Error(), res)
      } else {
        Send("OK", "Data purged!", res)
      }
  })

  mux.Put(prefix+"/data/load", func(res http.ResponseWriter, req *http.Request) {
      nodes := manager.FindBy(manager.SelectBuilder(), 0, 10)

      if nodes.Len() != 0 {
        Send("KO", "Table contains data, purge the data first!", res)

        return
      }

      err := LoadFixtures(manager, 100)

      if err != nil {
        Send("KO", err.Error(), res)
      } else {
        Send("OK", "Data loaded!", res)
      }
  })
}

func GetFakeMediaNode(pos int) *nc.Node {
  node := nc.NewNode()

  node.Type = "media.image"
  node.Name = "The image "+strconv.Itoa(pos)
  node.Slug = "the-image-"+strconv.Itoa(pos)
  node.Data = &nh.Image{
    Name: "Go pic",
    Reference: "0x0",
  }
  node.Meta = &nh.ImageMeta{}

  return node
}

func GetFakePostNode(pos int) *nc.Node {
  node := nc.NewNode()

  node.Type = "blog.post"
  node.Name = "The blog post "+strconv.Itoa(pos)
  node.Slug = "the-blog-post-"+strconv.Itoa(pos)
  node.Data = &nh.Post{
    Title: "Go pic",
    Content: "The Content of my blog post",
    Tags: []string{"sport", "tennis", "soccer"},
  }
  node.Meta = &nh.PostMeta{
    Format: "markdown",
  }

  return node
}

func GetFakeUserNode(pos int) *nc.Node {
  node := nc.NewNode()

  node.Type = "core.user"
  node.Name = "The user "+strconv.Itoa(pos)
  node.Slug = "the-user-"+strconv.Itoa(pos)
  node.Data = &nh.User{
    Login: "user"+strconv.Itoa(pos),
    NewPassword: "user"+strconv.Itoa(pos),
  }
  node.Meta = &nh.UserMeta{
    PasswordCost: 12,
    PasswordAlgo: "bcrypt",
  }

  return node
}

func LoadFixtures(m *nc.PgNodeManager, max int) error {

  var err error

  // create user
  admin := nc.NewNode()

  admin.Uuid = nc.GetRootReference()
  admin.Type = "core.user"
  admin.Name = "The admin user"
  admin.Slug = "the-admin-user"
  admin.Data = &nh.User{
    Login: "admin",
    NewPassword: "admin",
  }
  admin.Meta = &nh.UserMeta{
    PasswordCost: 12,
    PasswordAlgo: "bcrypt",
  }

  m.Save(admin)

  for i := 1; i < max; i++ {
    node := GetFakeUserNode(i)
    node.UpdatedBy = admin.Uuid
    node.CreatedBy = admin.Uuid

    m.Save(node)

    if err != nil {
      panic(err)
    }
  }

  for i := 1; i < max; i++ {
    node := GetFakeMediaNode(i)
    node.UpdatedBy = admin.Uuid
    node.CreatedBy = admin.Uuid

    m.Save(node)

    if err != nil {
      panic(err)
    }
  }

  for i := 1; i < max; i++ {
    node := GetFakePostNode(i)
    node.UpdatedBy = admin.Uuid
    node.CreatedBy = admin.Uuid

    m.Save(node)

    if err != nil {
      panic(err)
    }
  }

  return nil
}
