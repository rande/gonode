// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/rande/gonode/core"
	"io"
)

type YoutubeMeta struct {
	Type            string `json:"type"`
	Html            string `json:"html"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	ProviderName    string `json:"provider_name"`
	AuthorName      string `json:"author_name"`
	AuthorUrl       string `json:"author_url"`
	ProviderUrl     string `json:"provider_url"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailHeight int    `json:"thumbnail_height"`
}

type Youtube struct {
	Vid    string `json:"vid,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type YoutubeHandler struct {
}

func (h *YoutubeHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &Youtube{
		Status: core.ProcessStatusInit,
	}, &YoutubeMeta{}
}

func (h *YoutubeHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == core.ProcessStatusInit {
		node.Data.(*Youtube).Status = core.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == core.ProcessStatusInit {
		node.Data.(*Youtube).Status = core.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == core.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

func (h *YoutubeHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == core.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

func (h *YoutubeHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *YoutubeHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *YoutubeHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *YoutubeHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}

type YoutubeListener struct {
	HttpClient core.HttpClient
}

func (l *YoutubeListener) Handle(notification *pq.Notification, m core.NodeManager) (int, error) {
	reference, err := core.GetReferenceFromString(notification.Extra)

	if err != nil { // unable to parse the reference
		return core.PubSubListenContinue, nil
	}

	node := m.Find(reference)

	if node == nil {
		return core.PubSubListenContinue, nil
	}

	if node.Data.(*Youtube).Status == core.ProcessStatusDone {
		return core.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(fmt.Sprintf("https://www.youtube.com/oembed?url=http://www.youtube.com/watch?v=%s&format=json", node.Data.(*Youtube).Vid))
	if err != nil {
		node.Data.(*Youtube).Status = core.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while retrieving json response"
		m.Save(node, true)

		return core.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(node.Meta.(*YoutubeMeta))

	if err != nil {
		node.Data.(*Youtube).Status = core.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while decoding json"
		m.Save(node, true)

		return core.PubSubListenContinue, err
	}

	node.Data.(*Youtube).Status = core.ProcessStatusDone

	m.Save(node, true)

	if node.Meta.(*YoutubeMeta).ThumbnailUrl != "" {
		image := m.NewNode("media.image")
		image.Data.(*Image).SourceUrl = node.Meta.(*YoutubeMeta).ThumbnailUrl
		image.Source = node.Uuid
		image.CreatedBy = node.CreatedBy
		image.UpdatedBy = node.UpdatedBy

		m.Save(image, false)
	}

	return core.PubSubListenContinue, nil
}
