// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package feed

import (
	"errors"
	"fmt"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/plugins/search"
	"io"
)

var (
	InvalidFormat = errors.New("Invalid feed format")
)

type Feed struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Index       *search.Index `json:"index"`
}

type FeedHandler struct {
}

func (h *FeedHandler) GetStruct() (core.NodeData, core.NodeMeta) {

	return &Feed{
		Index: &search.Index{
			Deleted: search.NewParam(false, "="),
		},
		Title: "Feed",
	}, &search.IndexMeta{}
}

func (h *FeedHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *FeedHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *FeedHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *FeedHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *FeedHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {
}

func (h *FeedHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *FeedHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *FeedHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}

type FeedViewHandler struct {
	Search  *search.SearchPGSQL
	Manager core.NodeManager
}

func (v *FeedViewHandler) Execute(node *core.Node, request *core.ViewRequest, response *core.ViewResponse) error {
	feed := node.Data.(*Feed)

	// we just copy over node to create search form
	searchForm := search.NewSearchFormFromIndex(feed.Index)
	searchForm.PerPage = 32
	searchForm.Page = 1

	pager := search.GetPager(searchForm, v.Manager, v.Search)

	if request.Format == "rss" {
		response.HttpResponse.Header().Set("Content-Type", "application/rss+xml")
	} else if request.Format == "atom" {
		response.HttpResponse.Header().Set("Content-Type", "application/atom+xml")
	} else {
		return InvalidFormat
	}

	response.
		Set(200, fmt.Sprintf("nodes/%s.%s.tpl", node.Type, request.Format)).
		Add("pager", pager)

	return nil
}
