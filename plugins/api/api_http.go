// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/helper"
	"github.com/rande/gonode/plugins/user"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/bcrypt"
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
			query = query.Where(core.NewExprSlice(fmt.Sprintf("%s ??| array["+sq.Placeholders(len(value))+"]", name), value))
		}

		if len(value) == 1 {
			name = GetJsonQuery(field+"."+name, "->>")
			query = query.Where(sq.Expr(fmt.Sprintf("%s = ?", name), value[0]))
		}
	}

	return query
}

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Prepare(func(app *goapp.App) error {
		app.Set("gonode.websocket.clients", func(app *goapp.App) interface{} {
			return list.New()
		})

		sub := app.Get("gonode.postgres.subscriber").(*core.Subscriber)
		sub.ListenMessage(conf.Databases["master"].Prefix+"_manager_action", func(notification *pq.Notification) (int, error) {
			logger := app.Get("logger").(*log.Logger)
			logger.Printf("WebSocket: Sending message \n")
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			for e := webSocketList.Front(); e != nil; e = e.Next() {
				if err := e.Value.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(notification.Extra)); err != nil {
					logger.Printf("Error writing to websocket")
				}
			}

			logger.Printf("WebSocket: End Sending message \n")

			return core.PubSubListenContinue, nil
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
		app.Get("gonode.postgres.subscriber").(*core.Subscriber).Register()

		return nil
	})

	l.Exit(func(app *goapp.App) error {
		logger := app.Get("logger").(*log.Logger)
		logger.Printf("Closing PostgreSQL subcriber \n")
		app.Get("gonode.postgres.subscriber").(*core.Subscriber).Stop()
		logger.Printf("End closing PostgreSQL subcriber \n")

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		manager := app.Get("gonode.manager").(*core.PgNodeManager)
		apiHandler := app.Get("gonode.api").(*Api)
		handler_collection := app.Get("gonode.handler_collection").(core.Handlers)

		prefix := ""

		mux.Get(prefix+"/hello", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Hello!"))
		})

		mux.Post(prefix+"/login", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			req.ParseForm()

			loginForm := &struct {
				Username string `schema:"username"`
				Password string `schema:"password"`
			}{}

			decoder := schema.NewDecoder()
			err := decoder.Decode(loginForm, req.Form)

			core.PanicOnError(err)

			query := manager.SelectBuilder().Where("type = 'core.user' AND data->>'username' = ?", loginForm.Username)

			node := manager.FindOneBy(query)

			password := []byte("$2a$10$KDobsZdRDVnuMqvimYH82.Tnu3suk5xP7QzhQjlCo7Wy7d67xtYay")

			if node != nil {
				data := node.Data.(*user.User)
				password = []byte(data.Password)
			}

			fmt.Printf("%s => %s", password, loginForm.Password)

			if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginForm.Password)); err == nil { // equal
				token := jwt.New(jwt.SigningMethodHS256)
				token.Header["kid"] = "the sha1"

				// Set some claims
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
				// Sign and get the complete encoded token as a string
				tokenString, err := token.SignedString([]byte(conf.Guard.Key))

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, "Unable to sign the token")
					return
				}

				core.PanicOnError(err)
				res.Write([]byte(tokenString))
			} else {
				helper.SendWithHttpCode(res, http.StatusForbidden, "Unable to authenticate request: "+err.Error())
			}
		})

		mux.Get(prefix+"/nodes/stream", func(res http.ResponseWriter, req *http.Request) {
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			upgrader.CheckOrigin = func(r *http.Request) bool {
				return true
			}

			ws, err := upgrader.Upgrade(res, req, nil)

			core.PanicOnError(err)

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
			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // ask for binary content
				reference, err := core.GetReferenceFromString(c.URLParams["uuid"])

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

					return
				}

				node := manager.Find(reference)

				if node == nil {
					helper.SendWithHttpCode(res, http.StatusNotFound, "Element not found")

					return
				}

				data := handler_collection.Get(node).GetDownloadData(node)

				res.Header().Set("Content-Type", data.ContentType)

				//			if download {
				//				res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", data.Filename));
				//			}

				data.Stream(node, res)
			} else {
				// send the json value
				res.Header().Set("Content-Type", "application/json")
				err := apiHandler.FindOne(c.URLParams["uuid"], res)

				if err == core.NotFoundError {
					helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
				}

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
				}
			}
		})

		mux.Post(prefix+"/nodes", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			w := bufio.NewWriter(res)

			err := apiHandler.Save(req.Body, w)

			if err == core.RevisionError {
				res.WriteHeader(http.StatusConflict)
			}

			if err == core.ValidationError {
				res.WriteHeader(http.StatusPreconditionFailed)
			}

			res.WriteHeader(http.StatusCreated)

			w.Flush()
		})

		mux.Put(prefix+"/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // send binary data
				reference, err := core.GetReferenceFromString(c.URLParams["uuid"])

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

					return
				}

				node := manager.Find(reference)

				if node == nil {
					helper.SendWithHttpCode(res, http.StatusNotFound, "Element not found")
					return
				}

				_, err = handler_collection.Get(node).StoreStream(node, req.Body)

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
				} else {
					manager.Save(node, false)

					helper.SendWithHttpCode(res, http.StatusOK, "binary stored")
				}

			} else {
				w := bufio.NewWriter(res)

				err := apiHandler.Save(req.Body, w)

				if err == core.RevisionError {
					res.WriteHeader(http.StatusConflict)
				}

				if err == core.ValidationError {
					res.WriteHeader(http.StatusPreconditionFailed)
				}

				w.Flush()
			}
		})

		mux.Put(prefix+"/nodes/move/:uuid/:parentUuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			err := apiHandler.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], res)

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Delete(prefix+"/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			err := apiHandler.RemoveOne(c.URLParams["uuid"], res)

			if err == core.NotFoundError {
				helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
				return
			}

			if err == core.AlreadyDeletedError {
				helper.SendWithHttpCode(res, http.StatusGone, err.Error())
				return
			}

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Put(prefix+"/notify/:name", func(c web.C, res http.ResponseWriter, req *http.Request) {
			body, _ := ioutil.ReadAll(req.Body)

			manager.Notify(c.URLParams["name"], string(body[:]))
		})

		mux.Get(prefix+"/nodes", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			query := apiHandler.SelectBuilder()

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
				helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid pagination range")

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
					helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid order_by condition")

					return
				}

				query = query.OrderBy(GetJsonQuery(r[0][1], "->") + " " + r[0][2])
			}

			query = GetJsonSearchQuery(query, searchForm.Meta, "meta")
			query = GetJsonSearchQuery(query, searchForm.Data, "data")

			apiHandler.Find(res, query, uint64(searchForm.Page), uint64(searchForm.PerPage))
		})

		return nil
	})
}
