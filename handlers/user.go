package handlers

import (
	nc "github.com/rande/gonode/core"
)

type UserMeta struct {

}

type User struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserHandler struct {

}

func (h *UserHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &User{}, &UserMeta{}
}

func (h *UserHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}
