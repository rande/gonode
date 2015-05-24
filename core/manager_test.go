package core

import (
	"github.com/stretchr/testify/assert"
	"testing"

	//	"github.com/twinj/uuid"
)

func Test_Manager_Validate(t *testing.T) {
	m := &PgNodeManager{
		Handlers: map[string]Handler{
			"core.user": &UserHandler{},
		},
	}

	n := m.NewNode("core.user")

	ok, errors := m.Validate(n)

	assert.False(t, ok)
	assert.True(t, errors.HasErrors())
}

type UserMeta struct {
}

type User struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
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

	if data.Login == "" {
		errors.AddError("data.login", "Login cannot be empty")
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
