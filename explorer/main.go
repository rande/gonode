package main

import (
	nc "github.com/rande/gonode/core"
	nh "github.com/rande/gonode/handlers"
	"github.com/rande/gonode/extra"
	"net/http"
	pq "github.com/lib/pq"
	"database/sql"
	"log"
	sq "github.com/lann/squirrel"
	"github.com/spf13/afero"
	"strconv"
	"fmt"
	"time"
	"os"
	"github.com/hypebeast/gojistatic"
	"github.com/zenazn/goji"
)

func main() {
	manager := GetManager(nil)

//	LoadFixtures(manager, 256)
//	Check(manager, manager.Logger)


	extra.ConfigureGoji(manager, "")

	ListenMessages(manager)

  goji.Use(gojistatic.Static("dist", gojistatic.StaticOptions{SkipLogging: false, Prefix: "dist"}))

  // start http server
	goji.Serve()
}

func ListenMessage(name string, listenerHandler nc.Listener, manager nc.NodeManager) {
	var conninfo string = "postgres://safre:safre@localhost/safre"

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// listen to the specific channel
	listener := pq.NewListener(conninfo, 10 * time.Second, time.Minute, reportProblem)
	err := listener.Listen(name)
	if err != nil {
		panic(err)
	}

	// iterate over received notifications, for now, we start only one consumer with no concurrence
	for {
		select {
			case notification := <-listener.Notify:

				if notification == nil {
					fmt.Println("received a nil notification, the underlying driver reconnect")
					continue
				}
				fmt.Println("received notification, new work available")

				listenerHandler.Handle(notification, manager)

			case <-time.After(90 * time.Second):
				go func() {
					listener.Ping()
				}()
				// Check if there's more work available, just in case it takes
				// a while for the Listener to notice connection loss and
				// reconnect.
				fmt.Println("received no work for 90 seconds, checking for new work")
		}
	}
}

func ListenMessages(manager nc.NodeManager) {
	for name, listener := range GetListeners()  {
		go ListenMessage(name, listener, manager)
	}
}

func GetHandlers() map[string]nc.Handler {
	return map[string]nc.Handler{
		"default":       &nh.DefaultHandler{},
		"media.image":   &nh.ImageHandler{},
		"media.youtube": &nh.YoutubeHandler{},
		"blog.post":     &nh.PostHandler{},
		"core.user":     &nh.UserHandler{},
	}
}

func GetListeners() map[string]nc.Listener {
	return map[string]nc.Listener{
		"media_youtube_update": &nh.YoutubeListener{
			HttpClient: &http.Client{},
		},
		"media_file_download": &nh.ImageDownloadListener{
			Fs: nc.NewSecureFs(&afero.OsFs{}, "/tmp/gnode"),
			HttpClient: &http.Client{},
		},
	}
}

func GetManager(logger *log.Logger) *nc.PgNodeManager {

	if logger == nil {
		logger = log.New(os.Stdout, "Notification > ", log.Lshortfile)
	}

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

	return &nc.PgNodeManager{
		Logger: logger,
		Db: db,
		ReadOnly: false,
		Handlers: GetHandlers(),
		Listeners: GetListeners(),
	}
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
		Password: "{plain}user"+strconv.Itoa(pos),
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

	admin.Uuid = nc.GetRootReference()
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
