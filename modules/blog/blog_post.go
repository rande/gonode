// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package blog

import (
	"github.com/rande/gonode/modules/base"
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

func (h *PostHandler) GetMetadata() *base.HandlerMetadata {
	meta := base.NewHandlerMetadata()

	meta.Authors = []string{"Thomas Rabaix <thomas.rabaix@gmail.com>"}
	meta.Description = "Blog post engine"
	meta.Name = "Blog post"

	return meta
}
