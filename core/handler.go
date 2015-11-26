// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"io"
)

type NodeData interface{}
type NodeMeta interface{}

type Handlers interface {
	NewNode(t string) *Node
	Get(node *Node) Handler
	GetByCode(code string) Handler
	GetKeys() []string
}

type HandlerCollection map[string]Handler

func (c HandlerCollection) NewNode(t string) *Node {
	node := NewNode()
	node.Type = t
	node.Data, node.Meta = c.Get(node).GetStruct()

	return node
}

func (c HandlerCollection) Get(node *Node) Handler {
	return c.GetByCode(node.Type)
}

func (c HandlerCollection) GetByCode(code string) Handler {
	if handler, ok := c[code]; ok {
		return handler.(Handler)
	}

	return c["default"].(Handler)
}

func (c HandlerCollection) GetKeys() []string {
	keys := make([]string, 0)

	for k := range c {
		keys = append(keys, k)
	}

	return keys
}

type DownloadData struct {
	ContentType  string
	Filename     string
	CacheControl string
	Pragma       string
	Expires      string
	Stream       func(node *Node, w io.Writer)
}

type Handler interface {
	GetStruct() (NodeData, NodeMeta) // Data, Meta
	PreUpdate(node *Node, m NodeManager) error
	PostUpdate(node *Node, m NodeManager) error
	PreInsert(node *Node, m NodeManager) error
	PostInsert(node *Node, m NodeManager) error
	Validate(node *Node, m NodeManager, e Errors)
	Load(data []byte, meta []byte, node *Node) error
	GetDownloadData(node *Node) *DownloadData
	StoreStream(node *Node, r io.Reader) (int64, error)
}

func GetDownloadData() *DownloadData {
	return &DownloadData{
		ContentType:  "application/octet-stream",
		Filename:     "gonode-notype.bin",
		CacheControl: "private",
		Stream: func(node *Node, w io.Writer) {
			io.WriteString(w, "No content defined to be download for this node")
		},
	}
}

func HandlerLoad(handler Handler, data []byte, meta []byte, node *Node) error {
	var err error

	node.Data, node.Meta = handler.GetStruct()

	err = json.Unmarshal(data, node.Data)
	PanicOnError(err)

	err = json.Unmarshal(meta, node.Meta)
	PanicOnError(err)

	return nil
}

func DefaultHandlerStoreStream(node *Node, r io.Reader) (int64, error) {
	return 0, NoStreamHandler
}
