// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"github.com/lib/pq"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/modules/helper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_YoutubeHandler(t *testing.T) {
	a := assert.New(t)

	handler := &YoutubeHandler{}

	data, meta := handler.GetStruct()

	a.IsType(&YoutubeMeta{}, meta)
	a.IsType(&Youtube{}, data)

	a.Equal(data.(*Youtube).Status, core.ProcessStatusInit)
	a.Equal(data.(*Youtube).Vid, "")
}

func Test_YoutubeHandler_PreInsert(t *testing.T) {
	a := assert.New(t)

	node := core.NewNode()

	handler := &YoutubeHandler{}
	manager := &core.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	//no url => keep status
	handler.PreInsert(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusInit)

	//url => update status
	node.Data.(*Youtube).Vid = "k72S8XYqi0c"
	handler.PreInsert(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusUpdate)

	//url, status done => no update
	node.Data.(*Youtube).Status = core.ProcessStatusDone
	handler.PreInsert(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusDone)
}

func Test_YoutubeHandler_PreUpdate(t *testing.T) {
	a := assert.New(t)

	node := core.NewNode()

	handler := &YoutubeHandler{}
	manager := &core.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusInit)

	// url => update status
	node.Data.(*Youtube).Vid = "k72S8XYqi0c"
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusUpdate)

	// url, status done => no update
	node.Data.(*Youtube).Status = core.ProcessStatusDone
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusDone)
}

func Test_YoutubeHandler_PostUpdate(t *testing.T) {
	a := assert.New(t)

	node := core.NewNode()

	handler := &YoutubeHandler{}
	manager := &core.MockedManager{}
	manager.On("Notify", "media_youtube_update", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostUpdate(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusInit)

	// url => keep status
	node.Data.(*Youtube).Vid = "k72S8XYqi0c"
	node.Data.(*Youtube).Status = core.ProcessStatusUpdate

	handler.PostUpdate(node, manager)

	manager.AssertCalled(t, "Notify", "media_youtube_update", node.Uuid.String())

	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusUpdate)
}

func Test_YoutubeHandler_PostInsert(t *testing.T) {
	a := assert.New(t)

	node := core.NewNode()

	handler := &YoutubeHandler{}
	manager := &core.MockedManager{}
	manager.On("Notify", "media_youtube_update", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostInsert(node, manager)
	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusInit)

	// url => keep status
	node.Data.(*Youtube).Vid = "k72S8XYqi0c"
	node.Data.(*Youtube).Status = core.ProcessStatusUpdate

	handler.PostInsert(node, manager)

	manager.AssertCalled(t, "Notify", "media_youtube_update", node.Uuid.String())

	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusUpdate)
}

func Test_YoutubeListener_NodeNotFound(t *testing.T) {
	client := &helper.MockedHttpClient{}

	l := &YoutubeListener{
		HttpClient: client,
	}

	manager := &core.MockedManager{}
	manager.On("Find", core.GetEmptyReference()).Return(nil)

	notification := &pq.Notification{
		Extra: "11111111-1111-1111-1111-111111111111",
	}

	l.Handle(notification, manager)

	manager.AssertCalled(t, "Find", core.GetEmptyReference())
	manager.AssertNotCalled(t, "Save", nil)
}

func Test_YoutubeListener_Found(t *testing.T) {
	a := assert.New(t)

	handler := &YoutubeHandler{}
	node := core.NewNode()

	node.Data, node.Meta = handler.GetStruct()

	node.Data.(*Youtube).Status = core.ProcessStatusUpdate
	node.Data.(*Youtube).Vid = "MyVideoId"

	client := &helper.MockedHttpClient{}
	client.
		On("Get", "https://www.youtube.com/oembed?url=http://www.youtube.com/watch?v=MyVideoId&format=json").
		Return(&http.Response{Body: helper.NewTestCloserReader(`{
"provider_url": "http://www.youtube.com/",
"thumbnail_height": 360,
"thumbnail_url": "http://i.ytimg.com/vi/k72S8XYqi0c/hqdefault.jpg",
"type": "video",
"html": "<iframe width=\"480\" height=\"270\" src=\"http://www.youtube.com/embed/k72S8XYqi0c?feature=oembed\" frameborder=\"0\" allowfullscreen></iframe>",
"version": "1.0",
"author_name": "Comptines et chansons",
"height": 270,
"width": 480,
"provider_name": "YouTube",
"author_url": "http://www.youtube.com/user/comptines",
"thumbnail_width": 480,
"title": "La famille Tortue"
}`)}, nil)

	l := &YoutubeListener{
		HttpClient: client,
	}

	nodeImage := core.NewNode()
	nodeImage.Type = "media.image"
	nodeImage.Data = &Image{}
	nodeImage.Meta = &ImageMeta{}

	manager := &core.MockedManager{}
	manager.On("Find", core.GetEmptyReference()).Return(node)
	manager.On("Save", node).Return(node, nil)
	manager.On("Save", nodeImage).Return(nodeImage, nil)
	manager.On("NewNode", "media.image").Return(nodeImage, nil)

	notification := &pq.Notification{
		Extra: "11111111-1111-1111-1111-111111111111",
	}

	l.Handle(notification, manager)

	manager.AssertCalled(t, "Find", core.GetEmptyReference())
	manager.AssertCalled(t, "Save", node)
	manager.AssertCalled(t, "Save", nodeImage)
	client.AssertCalled(t, "Get", "https://www.youtube.com/oembed?url=http://www.youtube.com/watch?v=MyVideoId&format=json")

	a.Equal(nodeImage.Data.(*Image).SourceUrl, "http://i.ytimg.com/vi/k72S8XYqi0c/hqdefault.jpg")
	a.Equal(node.Meta.(*YoutubeMeta).ThumbnailUrl, "http://i.ytimg.com/vi/k72S8XYqi0c/hqdefault.jpg")
	a.Equal(node.Meta.(*YoutubeMeta).ProviderUrl, "http://www.youtube.com/")
	a.Equal(node.Meta.(*YoutubeMeta).ProviderName, "YouTube")
	a.Equal(node.Meta.(*YoutubeMeta).Type, "video")
	a.Equal(node.Meta.(*YoutubeMeta).Html, "<iframe width=\"480\" height=\"270\" src=\"http://www.youtube.com/embed/k72S8XYqi0c?feature=oembed\" frameborder=\"0\" allowfullscreen></iframe>")
	a.Equal(node.Meta.(*YoutubeMeta).ThumbnailHeight, 360)
	a.Equal(node.Meta.(*YoutubeMeta).ThumbnailWidth, 480)
	a.Equal(node.Meta.(*YoutubeMeta).Title, "La famille Tortue")
	a.Equal(node.Meta.(*YoutubeMeta).Height, 270)
	a.Equal(node.Meta.(*YoutubeMeta).Width, 480)

	a.Equal(node.Data.(*Youtube).Status, core.ProcessStatusDone)
	a.Equal(nodeImage.CreatedBy, node.CreatedBy)
	a.Equal(nodeImage.UpdatedBy, node.UpdatedBy)
	a.Equal(nodeImage.Source, node.Source)
}
