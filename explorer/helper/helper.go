package helper

import (
	"encoding/json"
	"fmt"
	"github.com/hypebeast/gojistatic"
	. "github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/extra"
	nh "github.com/rande/gonode/handlers"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

func BuildApp(app *App, config string) *App {
	app.Set("gonode.configuration", func(app *App) interface{} {
		return extra.GetConfiguration(config)
	})

	// configure main services
	app.Set("logger", func(app *App) interface{} {
		return log.New(os.Stdout, "", log.Lshortfile)
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

	// load the current application
	extra.ConfigureApp(app)
	extra.ConfigureGoji(app)

	ConfigureGoji(app)

	return app
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
		"status":  status,
		"message": message,
	})

	res.Write(data)
}

func ConfigureGoji(app *App) {

	mux := app.Get("goji.mux").(*web.Mux)
	manager := app.Get("gonode.manager").(*nc.PgNodeManager)
	configuration := app.Get("gonode.configuration").(*extra.Config)
	prefix := ""

	mux.Put(prefix+"/data/purge", func(res http.ResponseWriter, req *http.Request) {

		prefix := configuration.Databases["master"].Prefix

		tx, _ := manager.Db.Begin()
		manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes"`, prefix))
		manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes_audit"`, prefix))
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
	node.Name = "The image " + strconv.Itoa(pos)
	node.Slug = "the-image-" + strconv.Itoa(pos)
	node.Data = &nh.Image{
		Name:      "Go pic",
		Reference: "0x0",
	}
	node.Meta = &nh.ImageMeta{}

	return node
}

func GetFakePostNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "blog.post"
	node.Name = "The blog post " + strconv.Itoa(pos)
	node.Slug = "the-blog-post-" + strconv.Itoa(pos)
	node.Data = &nh.Post{
		Title:   "Go pic",
		Content: "The Content of my blog post",
		Tags:    []string{"sport", "tennis", "soccer"},
	}
	node.Meta = &nh.PostMeta{
		Format: "markdown",
	}

	return node
}

func GetFakeUserNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "core.user"
	node.Name = "The user " + strconv.Itoa(pos)
	node.Slug = "the-user-" + strconv.Itoa(pos)
	node.Data = &nh.User{
		Login:       "user" + strconv.Itoa(pos),
		NewPassword: "user" + strconv.Itoa(pos),
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
		Login:       "admin",
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
