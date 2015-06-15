package handlers

import (
	nc "github.com/rande/gonode/core"
)

type DefaultHandler struct {
}

func (h *DefaultHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return make(map[string]interface{}), make(map[string]interface{})
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

func (h *DefaultHandler) GetDownloadData(node *nc.Node) *nc.DownloadData {
	return nc.GetDownloadData()
}

func (h *DefaultHandler) Load(data []byte, meta []byte, node *nc.Node) error {
	return nc.HandlerLoad(h, data, meta, node)
}
