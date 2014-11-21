package main

import (
	nc "github.com/rande/gonode/core"
	nh "github.com/rande/gonode/handlers"
	"net/http"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"os"
	sq "github.com/lann/squirrel"
	"strconv"
)


func main() {
	var err error
	var db *sql.DB

	sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	db, err = sql.Open("postgres", "postgres://safre:safre@localhost/safre")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "> ", log.Lshortfile)

	logger.Printf("[PgNode] %s\n", "Loading fixtures")

	defer db.Close()

	manager := &nc.PgNodeManager{
		Logger: logger,
		Db: db,
		ReadOnly: false,
		Handlers: map[string] interface {}{
			"media.image": &nh.MediaHandler{},
			"blog.post": &nh.PostHandler{},
			"core.user": &nh.UserHandler{},
		},
	}

	api := &nc.Api{
		Manager: manager,
		Version: "1.0.0",
	}

//	LoadFixtures(manager, 256)

//	Check(manager, logger)

	ServeMartini(api, "")
}

func ServeMartini(api *nc.Api, prefix string) {

	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Get(prefix + "/node/:uuid", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		api.FindOne(params["uuid"], res)
	})


	//	m.Put(prefix+"/node", func(res http.ResponseWriter, req *http.Request) {
//		node := NewNode()
//
//		DecodeNode(node, req)
//
//		manager.Save(node)
//
//		EncodeNode(node, res)
//	})
//
//	m.Post(prefix+"/node/:uuid", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
//		node := manager.FindOne(bson.M{"uuid": params["uuid"]})
//
//		DecodeNode(node, req)
//
//		manager.Save(node)
//
//		EncodeNode(node, res)
//	})


//	// Search a node
//	m.Get(prefix+"/node", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
//		query := bson.M{}
//
//		var limit int
//		var offset int
//
//		req.ParseForm()
//
//		if len(req.FormValue("name")) > 0 {
//			query["name"] = bson.RegEx{req.FormValue("name"), "i"}
//		}
//
//		if len(req.FormValue("set")) > 0 {
//			query["set"] = bson.RegEx{req.FormValue("set"), "i"}
//		}
//
//		if len(req.FormValue("limit")) > 0 {
//			limit, _ = strconv.Atoi(req.FormValue("limit"))
//		}
//
//		if len(req.FormValue("offset")) > 0 {
//			offset, _ = strconv.Atoi(req.FormValue("offset"))
//		}
//
//		if limit < 1 {
//			limit = 25
//		}
//
//		if limit > 10 {
//			limit = 10
//		}
//
//		if _, ok := req.Form["types"]; ok {
//			query["type"] = bson.M{
//				"$in": req.Form["types"],
//			}
//		}
//
//		nodes := manager.Find(query, offset, limit)
//
//		query["offet"] = offset
//		query["limit"] = limit
//
//		EncodeNode(bson.M{
//			"results":   nodes,
//			"filters": query,
//		}, res)
//	})
//

	m.Run()
}

func Check(manager *nc.PgNodeManager, logger *log.Logger) {
	query := manager.SelectBuilder()
	nodes := manager.FindBy(query, 0, 10)

	if nodes.Len() == 0 {
		panic("Loading fixture failed ...")
	}

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

func GetFakeMediaNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "media.image"
	node.Name = "The image " + strconv.Itoa(pos)
	node.Slug = "the-image-" + strconv.Itoa(pos)
	node.Data =  &nh.Media{
		Name: "Go pic",
		Reference: "0x0",
	}
	node.Meta = &nh.MediaMeta{}

	return node
}

func GetFakePostNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "blog.post"
	node.Name = "The blog post " + strconv.Itoa(pos)
	node.Slug = "the-blog-post-" + strconv.Itoa(pos)
	node.Data = &nh.Post{
		Title: "Go pic",
		Content: "The Content of my blog post",
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
		Login: "user" + strconv.Itoa(pos),
		Password: "{plain}user" + strconv.Itoa(pos),
	}
	node.Meta = &nh.UserMeta{}

	return node
}

func LoadFixtures(m *nc.PgNodeManager, max int) {

	var err error

	m.Db.Query("DELETE FROM nodes")
	m.Db.Query("DELETE FROM nodes_audit")

	// create user
	admin := nc.NewNode()

	admin.Uuid = nc.GetRootUuid()
	admin.Type = "core.user"
	admin.Name = "The admin user"
	admin.Slug = "the-admin-user"
	admin.Data = &nh.User{
		Login: "admin",
		Password: "{plain}admin",
	}
	admin.Meta = &nh.UserMeta{}

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
}
