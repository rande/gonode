// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	rexOrderBy = regexp.MustCompile(`(^[a-z,_.A-Z]*),(DESC|ASC|desc|asc)$`)
	rexMeta    = regexp.MustCompile(`meta\.([a-zA-Z]*)`)
	rexData    = regexp.MustCompile(`data\.([a-zA-Z]*)`)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			return
		}
	}
}

func GetJsonQuery(left string, sep string) string {
	fields := strings.Split(left, ".")

	c := ""
	for p, f := range fields {
		if p == 0 {
			c += f
		} else {
			c += fmt.Sprintf(sep+"'%s'", f)
		}
	}

	return c
}

func GetJsonSearchQuery(query sq.SelectBuilder, data map[string][]string, field string) sq.SelectBuilder {
	//-- SELECT uuid, "data" #> '{tags,1}' as tags FROM nodes WHERE  "data" @> '{"tags": ["sport"]}'
	//-- SELECT uuid, "data" #> '{tags}' AS tags FROM nodes WHERE  "data" -> 'tags' ?| array['sport'];
	for name, value := range data {
		if len(value) > 1 {
			name = GetJsonQuery(field+"."+name, "->")
			query = query.Where(NewExprSlice(fmt.Sprintf("%s ??| array["+sq.Placeholders(len(value))+"]", name), value))
		}

		if len(value) == 1 {
			name = GetJsonQuery(field+"."+name, "->>")
			query = query.Where(sq.Expr(fmt.Sprintf("%s = ?", name), value[0]))
		}
	}

	return query
}

func ConfigureHttpApi(l *goapp.Lifecycle) {

	l.Prepare(func(app *goapp.App) error {
		app.Set("gonode.websocket.clients", func(app *goapp.App) interface{} {
			return list.New()
		})

		configuration := app.Get("gonode.configuration").(*ServerConfig)

		sub := app.Get("gonode.postgres.subscriber").(*Subscriber)
		sub.ListenMessage(configuration.Databases["master"].Prefix+"_manager_action", func(notification *pq.Notification) (int, error) {
			logger := app.Get("logger").(*log.Logger)
			logger.Printf("WebSocket: Sending message \n")
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			for e := webSocketList.Front(); e != nil; e = e.Next() {
				if err := e.Value.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(notification.Extra)); err != nil {
					logger.Printf("Error writing to websocket")
				}
			}

			logger.Printf("WebSocket: End Sending message \n")

			return PubSubListenContinue, nil
		})

		graceful.PreHook(func() {
			logger := app.Get("logger").(*log.Logger)
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			logger.Printf("Closing websocket connections \n")
			for e := webSocketList.Front(); e != nil; e = e.Next() {
				e.Value.(*websocket.Conn).Close()
			}
		})

		return nil
	})

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		logger := app.Get("logger").(*log.Logger)
		logger.Printf("Starting PostgreSQL subcriber \n")
		app.Get("gonode.postgres.subscriber").(*Subscriber).Register()

		return nil
	})

	l.Exit(func(app *goapp.App) error {
		logger := app.Get("logger").(*log.Logger)
		logger.Printf("Closing PostgreSQL subcriber \n")
		app.Get("gonode.postgres.subscriber").(*Subscriber).Stop()
		logger.Printf("End closing PostgreSQL subcriber \n")

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		manager := app.Get("gonode.manager").(*PgNodeManager)
		api := app.Get("gonode.api").(*Api)

		prefix := ""

		handlers := app.Get("gonode.handler_collection").(Handlers)
		configuration := app.Get("gonode.configuration").(*ServerConfig)

		mux.Get(prefix+"/hello", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)

			res.Write([]byte("Hello!"))
		})

		mux.Get(prefix+"/nodes/stream", func(res http.ResponseWriter, req *http.Request) {
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			upgrader.CheckOrigin = func(r *http.Request) bool {
				return true
			}

			ws, err := upgrader.Upgrade(res, req, nil)

			PanicOnError(err)

			element := webSocketList.PushBack(ws)

			var closeDefer = func() {
				ws.Close()
				webSocketList.Remove(element)
			}

			defer closeDefer()

			go readLoop(ws)

			// ping remote client, avoid keeping open connection
			for {
				time.Sleep(2 * time.Second)
				if err := ws.WriteMessage(websocket.TextMessage, []byte("PING")); err != nil {
					return
				}
			}
		})

		mux.Get(prefix+"/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // ask for binary content
				reference, err := GetReferenceFromString(c.URLParams["uuid"])

				if err != nil {
					SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

					return
				}

				node := manager.Find(reference)

				if node == nil {
					SendWithHttpCode(res, http.StatusNotFound, "Element not found")

					return
				}

				data := handlers.Get(node).GetDownloadData(node)

				res.Header().Set("Content-Type", data.ContentType)

				//			if download {
				//				res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", data.Filename));
				//			}

				data.Stream(node, res)
			} else {
				// send the json value
				res.Header().Set("Content-Type", "application/json")
				err := api.FindOne(c.URLParams["uuid"], res)

				if err == NotFoundError {
					SendWithHttpCode(res, http.StatusNotFound, err.Error())
				}

				if err != nil {
					SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
				}
			}
		})

		mux.Post(prefix+"/nodes", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			w := bufio.NewWriter(res)

			err := api.Save(req.Body, w)

			if err == RevisionError {
				res.WriteHeader(http.StatusConflict)
			}

			if err == ValidationError {
				res.WriteHeader(http.StatusPreconditionFailed)
			}

			res.WriteHeader(http.StatusCreated)

			w.Flush()
		})

		mux.Put(prefix+"/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // send binary data
				reference, err := GetReferenceFromString(c.URLParams["uuid"])

				if err != nil {
					SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

					return
				}

				node := manager.Find(reference)

				if node == nil {
					SendWithHttpCode(res, http.StatusNotFound, "Element not found")
					return
				}

				_, err = handlers.Get(node).StoreStream(node, req.Body)

				if err != nil {
					SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
				} else {
					manager.Save(node, false)

					SendWithHttpCode(res, http.StatusOK, "binary stored")
				}

			} else {
				w := bufio.NewWriter(res)

				err := api.Save(req.Body, w)

				if err == RevisionError {
					res.WriteHeader(http.StatusConflict)
				}

				if err == ValidationError {
					res.WriteHeader(http.StatusPreconditionFailed)
				}

				w.Flush()
			}
		})

		mux.Put(prefix+"/nodes/moves/:uuid/:parentUuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			err := api.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], res)

			if err != nil {
				SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Delete(prefix+"/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			err := api.RemoveOne(c.URLParams["uuid"], res)

			if err == NotFoundError {
				SendWithHttpCode(res, http.StatusNotFound, err.Error())
				return
			}

			if err == AlreadyDeletedError {
				SendWithHttpCode(res, http.StatusGone, err.Error())
				return
			}

			if err != nil {
				SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Put(prefix+"/notify/:name", func(c web.C, res http.ResponseWriter, req *http.Request) {
			body, _ := ioutil.ReadAll(req.Body)

			manager.Notify(c.URLParams["name"], string(body[:]))
		})

		mux.Put(prefix+"/uninstall", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			prefix := configuration.Databases["master"].Prefix

			manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes_audit"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_idx"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_current_idx"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_id_seq" CASCADE`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_audit_id_seq" CASCADE`, prefix))

			SendWithHttpCode(res, http.StatusOK, "Successfully delete tables!")
		})

		mux.Put(prefix+"/install", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			prefix := configuration.Databases["master"].Prefix
			tx, _ := manager.Db.Begin()

			// Create my table
			tx.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			tx.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes" (
				"id" INTEGER DEFAULT nextval('%s_nodes_id_seq'::regclass) NOT NULL UNIQUE,
				"uuid" UUid NOT NULL,
				"type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL,
				"name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL,
				"enabled" BOOLEAN DEFAULT 'true' NOT NULL,
				"current" BOOLEAN DEFAULT 'false' NOT NULL,
				"revision" INTEGER DEFAULT '1' NOT NULL,
				"version" INTEGER DEFAULT '1' NOT NULL,
				"status" INTEGER DEFAULT '0' NOT NULL,
				"deleted" BOOLEAN DEFAULT 'false' NOT NULL,
				"data" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"meta" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
				"source" UUid,
				"set_uuid" UUid,
				"parent_uuid" UUid,
				"parents" UUid[],
				"created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"created_by" UUid NOT NULL,
				"updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"updated_by" UUid NOT NULL,
				"weight" INTEGER DEFAULT '0' NOT NULL,
				PRIMARY KEY ( "id" ),
				CONSTRAINT "%s_slug" UNIQUE( "parent_uuid","slug","revision" ),
				CONSTRAINT "%s_uuid" UNIQUE( "revision","uuid" )
			)`, prefix, prefix, prefix, prefix))

			tx.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST )`, prefix, prefix))
			tx.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_current_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST, "current" ASC NULLS LAST )`, prefix, prefix))

			// Create Index
			tx.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_audit_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			tx.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes_audit" (
				"id" INTEGER DEFAULT nextval('%s_nodes_id_seq'::regclass) NOT NULL UNIQUE,
				"uuid" UUid NOT NULL,
				"type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL,
				"name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL,
				"enabled" BOOLEAN DEFAULT 'true' NOT NULL,
				"current" BOOLEAN DEFAULT 'false' NOT NULL,
				"revision" INTEGER DEFAULT '1' NOT NULL,
				"version" INTEGER DEFAULT '1' NOT NULL,
				"status" INTEGER DEFAULT '0' NOT NULL,
				"deleted" BOOLEAN DEFAULT 'false' NOT NULL,
				"data" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"meta" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
				"source" UUid,
				"set_uuid" UUid,
				"parent_uuid" UUid,
				"parents" UUid[],
				"created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"created_by" UUid NOT NULL,
				"updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"updated_by" UUid NOT NULL,
				"weight" INTEGER DEFAULT '0' NOT NULL,
				PRIMARY KEY ( "id" )
			)`, prefix, prefix))

			err := tx.Commit()

			if err != nil {
				SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			} else {
				SendWithHttpCode(res, http.StatusOK, "Successfully create tables!")
			}
		})

		mux.Get(prefix+"/nodes", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("X-Generator", "gonode - thomas.rabaix@gmail.com - v"+api.Version)
			res.Header().Set("Access-Control-Allow-Origin", "*")

			query := api.SelectBuilder()

			req.ParseForm()

			searchForm := GetSearchForm()
			decoder := schema.NewDecoder()
			decoder.Decode(searchForm, req.Form)

			// analyse Meta
			for name, value := range req.Form {
				values := rexMeta.FindStringSubmatch(name)

				if len(values) == 2 {
					searchForm.Meta[values[1]] = value
				}
			}

			// analyse Data
			for name, value := range req.Form {
				values := rexData.FindStringSubmatch(name)

				if len(values) == 2 {
					searchForm.Data[values[1]] = value
				}
			}

			if searchForm.Page < 0 || searchForm.PerPage < 0 || searchForm.PerPage > 128 {
				SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid pagination range")

				return
			}

			if searchForm.Page == 0 {
				searchForm.Page = 1
			}

			if searchForm.PerPage == 0 {
				searchForm.PerPage = 32
			}

			if len(searchForm.Source) != 0 {
				query = query.Where(sq.Eq{"source": searchForm.Source})
			}

			if searchForm.Enabled != "" {
				query = query.Where("enabled = ?", searchForm.Enabled)
			}

			if len(searchForm.Type) != 0 {
				query = query.Where(sq.Eq{"type": searchForm.Type})
			}

			if searchForm.Current != "" {
				query = query.Where("current = ?", searchForm.Current)
			}

			if searchForm.Deleted != "" { // TODO: only admin token can view deleted node
				query = query.Where("deleted = ?", searchForm.Deleted)
			} else {
				query = query.Where("deleted = ?", "f")
			}

			if len(searchForm.Uuid) != 0 {
				query = query.Where(sq.Eq{"uuid": searchForm.Uuid})
			}

			if len(searchForm.ParentUuid) != 0 {
				query = query.Where(sq.Eq{"parent_uuid": searchForm.ParentUuid})
			}

			if searchForm.Slug != "" {
				query = query.Where("slug = ?", searchForm.Slug)
			}

			if searchForm.Revision != "" {
				query = query.Where("revision = ?", searchForm.Revision)
			}

			if len(searchForm.Status) != 0 {
				query = query.Where(sq.Eq{"status": searchForm.Status})
			}

			for _, order := range searchForm.OrderBy {
				r := rexOrderBy.FindAllStringSubmatch(order, -1)

				if r == nil {
					SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid order_by condition")

					return
				}

				query = query.OrderBy(GetJsonQuery(r[0][1], "->") + " " + r[0][2])
			}

			query = GetJsonSearchQuery(query, searchForm.Meta, "meta")
			query = GetJsonSearchQuery(query, searchForm.Data, "data")

			api.Find(res, query, uint64(searchForm.Page), uint64(searchForm.PerPage))
		})

		return nil
	})
}
