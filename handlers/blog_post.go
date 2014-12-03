package handlers

import (
	nc "github.com/rande/gonode/core"
)

type PostMeta struct {
	Format string     `json:"format"`
}

type Post struct {
	Title    string   `json:"title"`
	SubTitle string   `json:"sub_title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
}

type PostHandler struct {

}

func (h *PostHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &Post{}, &PostMeta{}
}

func (h *PostHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}
