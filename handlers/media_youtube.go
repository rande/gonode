package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	nc "github.com/rande/gonode/core"
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

func (h *YoutubeHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &Youtube{
		Status: nc.ProcessStatusInit,
	}, &YoutubeMeta{}
}

func (h *YoutubeHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == nc.ProcessStatusInit {
		node.Data.(*Youtube).Status = nc.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == nc.ProcessStatusInit {
		node.Data.(*Youtube).Status = nc.ProcessStatusUpdate
	}

	return nil
}

func (h *YoutubeHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == nc.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

func (h *YoutubeHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	if node.Data.(*Youtube).Vid != "" && node.Data.(*Youtube).Status == nc.ProcessStatusUpdate {
		m.Notify("media_youtube_update", node.Uuid.String())
	}

	return nil
}

func (h *YoutubeHandler) Validate(node *nc.Node, m nc.NodeManager, errors nc.Errors) {

}

func (h *YoutubeHandler) GetDownloadData(node *nc.Node) *nc.DownloadData {
	return nc.GetDownloadData()
}

type YoutubeListener struct {
	HttpClient nc.HttpClient
}

func (l *YoutubeListener) Handle(notification *pq.Notification, m nc.NodeManager) (int, error) {
	reference := nc.GetReferenceFromString(notification.Extra)

	node := m.Find(reference)

	if node == nil {
		return nc.PubSubListenContinue, nil
	}

	if node.Data.(*Youtube).Status == nc.ProcessStatusDone {
		return nc.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(fmt.Sprintf("https://www.youtube.com/oembed?url=http://www.youtube.com/watch?v=%s&format=json", node.Data.(*Youtube).Vid))
	if err != nil {
		node.Data.(*Youtube).Status = nc.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while retrieving json response"
		m.Save(node)

		return nc.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(node.Meta.(*YoutubeMeta))

	if err != nil {
		node.Data.(*Youtube).Status = nc.ProcessStatusError
		node.Data.(*Youtube).Error = "Error while decoding json"
		m.Save(node)

		return nc.PubSubListenContinue, err
	}

	node.Data.(*Youtube).Status = nc.ProcessStatusDone

	m.Save(node)

	if node.Meta.(*YoutubeMeta).ThumbnailUrl != "" {
		image := m.NewNode("media.image")
		image.Data.(*Image).SourceUrl = node.Meta.(*YoutubeMeta).ThumbnailUrl
		image.Source = node.Uuid
		image.CreatedBy = node.CreatedBy
		image.UpdatedBy = node.UpdatedBy

		m.Save(image)
	}

	return nc.PubSubListenContinue, nil
}
