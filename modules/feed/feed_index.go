// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package feed

import (
	"errors"
	"fmt"

	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
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

func (h *FeedHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Feed{
		Index: &search.Index{
			Deleted: search.NewParam(false, "="),
		},
		Title: "Feed",
	}, &search.IndexMeta{}
}

type FeedViewHandler struct {
	Search  *search.SearchPGSQL
	Manager base.NodeManager
}

func (v *FeedViewHandler) Support(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) bool {
	return request.Format == "atom" || request.Format == "rss"
}

func (v *FeedViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	feed := node.Data.(*Feed)

	// we just copy over node to create search form
	searchForm := search.NewSearchFormFromIndex(feed.Index)
	searchForm.PerPage = 32
	searchForm.Page = 1

	// apply security access
	options := base.NewAccessOptionsFromToken(security.GetTokenFromContext(request.Context))
	pager := search.GetPager(searchForm, v.Manager, v.Search, options)

	if request.Format == "rss" {
		response.HttpResponse.Header().Set("Content-Type", "application/rss+xml")
	} else if request.Format == "atom" {
		response.HttpResponse.Header().Set("Content-Type", "application/atom+xml")
	} else {
		return InvalidFormat
	}

	response.
		Set(200, fmt.Sprintf("feed:nodes/%s.%s.tpl", node.Type, request.Format)).
		Add("pager", pager)

	return nil
}
