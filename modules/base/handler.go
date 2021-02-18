// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/rande/gonode/core/helper"
	"github.com/zenazn/goji/web"
)

type NodeData interface{}
type NodeMeta interface{}

type HandlerMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Code        string   `json:"code"`
	Require     []string `json:"require"`
	Authors     []string `json:"authors"`
}

func NewHandlerMetadata() *HandlerMetadata {
	return &HandlerMetadata{
		Icon: `<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8">
  <path d="M3 0l-.34.34c-.11.11-2.66 2.69-2.66 4.88 0 1.65 1.35 3 3 3s3-1.35 3-3c0-2.18-2.55-4.77-2.66-4.88l-.34-.34zm-1.5 4.72c.28 0 .5.22.5.5 0 .55.45 1 1 1 .28 0 .5.22.5.5s-.22.5-.5.5c-1.1 0-2-.9-2-2 0-.28.22-.5.5-.5z" transform="translate(1)" />
</svg>`,
	}
}

type HandlerViewMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Code        string   `json:"code"`
	Require     []string `json:"require"`
	Authors     []string `json:"authors"`
}

func NewViewHandlerMetadata() *HandlerViewMetadata {
	return &HandlerViewMetadata{
		Icon: `<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8">
  <path d="M3 0l-.34.34c-.11.11-2.66 2.69-2.66 4.88 0 1.65 1.35 3 3 3s3-1.35 3-3c0-2.18-2.55-4.77-2.66-4.88l-.34-.34zm-1.5 4.72c.28 0 .5.22.5.5 0 .55.45 1 1 1 .28 0 .5.22.5.5s-.22.5-.5.5c-1.1 0-2-.9-2-2 0-.28.22-.5.5-.5z" transform="translate(1)" />
</svg>`,
	}
}

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

func (c HandlerCollection) Add(code string, h Handler) {
	c[code] = h
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
}

type MetadataHandler interface {
	GetMetadata() *HandlerMetadata
}

type ViewMetadataHandler interface {
	GetViewMetadata() *HandlerViewMetadata
}

type DatabaseNodeHandler interface {
	PreUpdate(node *Node, m NodeManager) error
	PostUpdate(node *Node, m NodeManager) error
	PreInsert(node *Node, m NodeManager) error
	PostInsert(node *Node, m NodeManager) error
}

type ValidateNodeHandler interface {
	Validate(node *Node, m NodeManager, e Errors)
}

type LoadNodeHandler interface {
	Load(data []byte, meta []byte, node *Node) error
}

type DownloadNodeHandler interface {
	GetDownloadData(node *Node) *DownloadData
}

type StoreStreamNodeHandler interface {
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
	helper.PanicOnError(err)

	err = json.Unmarshal(meta, node.Meta)
	helper.PanicOnError(err)

	return nil
}

func DefaultHandlerStoreStream(node *Node, r io.Reader) (int64, error) {
	return 0, NoStreamHandler
}

type ViewHandlerCollection map[string]ViewHandler

func (c ViewHandlerCollection) Get(node *Node) ViewHandler {
	return c.GetByCode(node.Type)
}

func (c ViewHandlerCollection) GetByCode(code string) ViewHandler {
	if handler, ok := c[code]; ok {
		return handler.(ViewHandler)
	}

	return c["default"].(ViewHandler)
}

func (c ViewHandlerCollection) GetKeys() []string {
	keys := make([]string, 0)

	for k := range c {
		keys = append(keys, k)
	}

	return keys
}

func (c ViewHandlerCollection) Add(code string, h ViewHandler) {
	c[code] = h
}

type ViewHandler interface {
	Support(node *Node, request *ViewRequest, response *ViewResponse) bool
	Execute(node *Node, request *ViewRequest, response *ViewResponse) error
}

type ViewRequest struct {
	Format      string
	HttpRequest *http.Request
	Context     web.C
}

func NewViewResponse(res http.ResponseWriter) *ViewResponse {
	return &ViewResponse{
		StatusCode:   200,
		Context:      pongo2.Context{},
		HttpResponse: res,
	}
}

type ViewResponse struct {
	StatusCode   int
	Template     string
	Context      pongo2.Context
	HttpResponse http.ResponseWriter
}

func (r *ViewResponse) Set(code int, template string) *ViewResponse {
	r.StatusCode = code
	r.Template = template

	return r
}

func (r *ViewResponse) Add(name string, v interface{}) *ViewResponse {
	r.Context[name] = v

	return r
}
