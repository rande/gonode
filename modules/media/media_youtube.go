// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/rande/gonode/modules/base"
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

func (h *YoutubeHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Youtube{
		Status: base.ProcessStatusInit,
	}, &YoutubeMeta{}
}

func (h *YoutubeHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == base.ProcessStatusInit {
		node.Data.(*Youtube).Status = base.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == base.ProcessStatusInit {
		node.Data.(*Youtube).Status = base.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == base.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

func (h *YoutubeHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == base.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

type YoutubeListener struct {
	HttpClient base.HttpClient
}

func (l *YoutubeListener) Handle(notification *pq.Notification, m base.NodeManager) (int, error) {
	reference, err := base.GetReferenceFromString(notification.Extra)

	if err != nil { // unable to parse the reference
		return base.PubSubListenContinue, nil
	}

	node := m.Find(reference)

	if node == nil {
		return base.PubSubListenContinue, nil
	}

	if node.Data.(*Youtube).Status == base.ProcessStatusDone {
		return base.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(fmt.Sprintf("https://www.youtube.com/oembed?url=http://www.youtube.com/watch?v=%s&format=json", node.Data.(*Youtube).Vid))
	if err != nil {
		node.Data.(*Youtube).Status = base.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while retrieving json response"
		m.Save(node, true)

		return base.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(node.Meta.(*YoutubeMeta))

	if err != nil {
		node.Data.(*Youtube).Status = base.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while decoding json"
		m.Save(node, true)

		return base.PubSubListenContinue, err
	}

	node.Data.(*Youtube).Status = base.ProcessStatusDone

	m.Save(node, true)

	if node.Meta.(*YoutubeMeta).ThumbnailUrl != "" {
		image := m.NewNode("media.image")
		image.Data.(*Image).SourceUrl = node.Meta.(*YoutubeMeta).ThumbnailUrl
		image.Source = node.Uuid
		image.CreatedBy = node.CreatedBy
		image.UpdatedBy = node.UpdatedBy
		image.Access = node.Access

		m.Save(image, false)
	}

	return base.PubSubListenContinue, nil
}
