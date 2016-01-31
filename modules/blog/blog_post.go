// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package blog

import (
	"github.com/rande/gonode/modules/base"
	"io"
	"time"
)

type PostMeta struct {
	Format string `json:"format"`
}

type Post struct {
	Title           string    `json:"title"`
	SubTitle        string    `json:"sub_title"`
	Content         string    `json:"content"`
	PublicationDate time.Time `json:"publication_date"`
	Tags            []string  `json:"tags"`
}

type PostHandler struct {
}

func (h *PostHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Post{
		PublicationDate: time.Now(),
	}, &PostMeta{}
}

func (h *PostHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *PostHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *PostHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *PostHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *PostHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {

}

func (h *PostHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	return base.GetDownloadData()
}

func (h *PostHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *PostHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return base.DefaultHandlerStoreStream(node, r)
}
