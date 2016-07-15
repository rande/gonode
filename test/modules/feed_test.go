// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"encoding/xml"

	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/feed"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type RssFeedXml struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *feeds.RssFeed
}

func Setup_Feed_Data(app *goapp.App) *base.Node {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)

	home := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.root")
	home.Name = "Blog"
	home.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

	manager.Save(home, false)

	post := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
	post.Name = "Article 1"
	post.Slug = "article-1"
	post.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

	manager.Save(post, false)
	manager.Move(post.Uuid, home.Uuid)

	post2 := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
	post2.Name = "Article 2"
	post2.Slug = "article-2"
	post2.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

	manager.Save(post2, false)
	manager.Move(post2.Uuid, home.Uuid)

	index := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("feed.index")
	index.Name = "Feed"
	index.Slug = "feed"
	index.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

	data := index.Data.(*feed.Feed)
	data.Description = "Blog post list"
	data.Title = "Blog post title"
	data.Index.Enabled = search.NewParam(true)
	data.Index.Type = search.NewParam("blog.post")

	manager.Save(index, false)
	manager.Move(index.Uuid, home.Uuid)

	return index
}

func Test_Feed_RSS(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		index := Setup_Feed_Data(app)

		// WHEN
		fp := gofeed.NewParser()
		f, err := fp.ParseURL(fmt.Sprintf("%s/prism/%s.rss", ts.URL, index.Uuid))

		assert.NoError(t, err)

		// THEN
		assert.Equal(t, "Blog post title", f.Title)
		assert.Equal(t, 2, len(f.Items))
	})
}

func Test_Feed_Atom(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		index := Setup_Feed_Data(app)

		// WHEN
		fp := gofeed.NewParser()
		f, err := fp.ParseURL(fmt.Sprintf("%s/prism/%s.atom", ts.URL, index.Uuid))

		assert.NoError(t, err)

		// THEN
		assert.Equal(t, "Blog post title", f.Title)
		assert.Equal(t, 2, len(f.Items))
	})
}