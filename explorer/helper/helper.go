// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"encoding/json"
	nc "github.com/rande/gonode/core"
	nh "github.com/rande/gonode/handlers"
	"log"
	"net/http"
	"strconv"
)

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

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakeMediaNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakePostNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	return nil
}
