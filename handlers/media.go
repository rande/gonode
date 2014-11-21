package handlers

import (
	gn "github.com/rande/gonode/core"
)

type MediaMeta struct {
	Width       int
	Height      int
	Size        int
	ContentType int
	Length      int
}

type Media struct {
	Reference string
	Name string
}

type MediaHandler struct {

}

func (h *MediaHandler) GetStruct() (gn.NodeData, gn.NodeMeta) {
	return &Media{}, &MediaMeta{}
}
