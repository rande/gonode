// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Manager_Validate(t *testing.T) {
	c := HandlerCollection{
		"node.user": &UserHandler{},
	}

	m := &PgNodeManager{
		Handlers: c,
	}

	n := c.NewNode("node.user")

	ok, errors := m.Validate(n)

	assert.False(t, ok)
	assert.True(t, errors.HasErrors())
}

type UserMeta struct {
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
}

func (h *UserHandler) GetStruct() (NodeData, NodeMeta) {
	return &User{}, &UserMeta{}
}

func (h *UserHandler) PreInsert(node *Node, m NodeManager) error {
	return nil
}

func (h *UserHandler) PreUpdate(node *Node, m NodeManager) error {
	return nil
}

func (h *UserHandler) PostInsert(node *Node, m NodeManager) error {
	return nil
}

func (h *UserHandler) PostUpdate(node *Node, m NodeManager) error {
	return nil
}

func (h *UserHandler) Validate(node *Node, m NodeManager, errors Errors) {

	data := node.Data.(*User)

	if data.Username == "" {
		errors.AddError("data.username", "Username cannot be empty")
	}

	if data.Name == "" {
		errors.AddError("data.name", "Name cannot be empty")
	}

	if data.Password == "" {
		errors.AddError("data.password", "Password cannot be empty")
	}
}

func (h *UserHandler) GetDownloadData(node *Node) *DownloadData {
	return GetDownloadData()
}

func (h *UserHandler) Load(data []byte, meta []byte, node *Node) error {
	return HandlerLoad(h, data, meta, node)
}

func (h *UserHandler) StoreStream(node *Node, r io.Reader) (int64, error) {
	return DefaultHandlerStoreStream(node, r)
}
