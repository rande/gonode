package handlers

import (
	nc "github.com/rande/gonode/core"
)

type PostMeta struct {
	Format string
}

type Post struct {
	Title    string
	SubTitle string
	Content  string
	Tags     []string
}

type PostHandler struct {

}

func (h *PostHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &Post{}, &PostMeta{}
}
