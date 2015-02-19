package handlers

import (
	nc "github.com/rande/gonode/core"
	"github.com/lib/pq"
	"github.com/twinj/uuid"
	"fmt"
	"github.com/spf13/afero"
	"io"
)

type ExifMeta map[string]string

type ImageMeta struct {
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Size        int      `json:"size"`
	ContentType int      `json:"content_type"`
	Length      int      `json:"length"`
	Exif        ExifMeta `json:"exif"`
	Hash        string   `json:"hash"`
}

type Image struct {
	Reference    string   `json:"reference"`
	Name         string   `json:"name"`
	SourceUrl    string   `json:"source_url"`
	SourceStatus int      `json:"source_status"`
	SourceError  string   `json:"source_error"`
}

type ImageHandler struct {

}

func (h *ImageHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &Image{
		SourceStatus: nc.ProcessStatusInit,
	}, &ImageMeta{}
}

func (h *ImageHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	image := node.Data.(*Image)

	if image.SourceUrl != "" && image.SourceStatus == nc.ProcessStatusInit {
		image.SourceStatus = nc.ProcessStatusUpdate
	}

	return nil
}

func (h *ImageHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	image := node.Data.(*Image)

	if image.SourceUrl != "" && image.SourceStatus == nc.ProcessStatusInit {
		image.SourceStatus = nc.ProcessStatusUpdate
	}

	return nil
}

func (h *ImageHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	image := node.Data.(*Image)

	if image.SourceStatus == nc.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	image := node.Data.(*Image)

	if image.SourceStatus == nc.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) Validate(node *nc.Node, m nc.NodeManager, errors nc.Errors) {

}

type ImageDownloadListener struct {
	Fs afero.Fs
	HttpClient nc.HttpClient
}

func (l *ImageDownloadListener) Handle(notification *pq.Notification, m nc.NodeManager) {

	reference := nc.GetReferenceFromString(notification.Extra)

	fmt.Printf("Download binary from uuid: %s\n", notification.Extra)
	node := m.Find(reference)

	if node == nil {
		fmt.Printf("Uuid does not exist: %s\n", notification.Extra)
		return
	}

	image := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if image.SourceStatus == nc.ProcessStatusDone {
		fmt.Printf("Nothing to update: %s\n", notification.Extra)

		return
	}

	resp, err := l.HttpClient.Get(image.SourceUrl)

	if err != nil {
		image.SourceStatus = nc.ProcessStatusError
		image.SourceError = "Unable to retrieve the remote file"
		m.Save(node)

		panic(err)
	}

	defer resp.Body.Close()

	strUuid := uuid.Formatter(node.Uuid, uuid.CleanHyphen)

	l.Fs.MkdirAll(fmt.Sprintf("%s/%s", strUuid[0:2], strUuid[2:4]), 0755)

	file, _ := l.Fs.Create(fmt.Sprintf("%s/%s/%s.bin", strUuid[0:2], strUuid[2:4], strUuid[4:]))

	written, err := io.Copy(file, resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Binary length: %ds\n", written)

	meta.Size = int(written)
	image.SourceStatus = nc.ProcessStatusDone

	m.Save(node)
}
