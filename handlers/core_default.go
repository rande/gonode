package handlers

import (
	nc "github.com/rande/gonode/core"
	"encoding/json"
)

type DefaultHandler struct {

}

func (h *DefaultHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &json.RawMessage{}, &json.RawMessage{}
}

func (h *DefaultHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *DefaultHandler) Validate(node *nc.Node, m nc.NodeManager, errors nc.Errors) {

}
