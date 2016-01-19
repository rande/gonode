// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/modules/vault"
	"io"
)

type ExifMeta map[string]string

type ImageMeta struct {
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Size         int      `json:"size"`
	ContentType  string   `json:"content_type"`
	Length       int      `json:"length"`
	Exif         ExifMeta `json:"exif"`
	Hash         string   `json:"hash"`
	SourceStatus int      `json:"source_status"`
	SourceError  string   `json:"source_error"`
}

type Image struct {
	Reference string `json:"reference"`
	Name      string `json:"name"`
	SourceUrl string `json:"source_url"`
}

type ImageHandler struct {
	Vault *vault.Vault
}

func (h *ImageHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &Image{}, &ImageMeta{
		SourceStatus: core.ProcessStatusInit,
	}
}

func (h *ImageHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if data.SourceUrl != "" && meta.SourceStatus == core.ProcessStatusInit {
		meta.SourceStatus = core.ProcessStatusUpdate
		meta.SourceError = ""
	}

	return nil
}

func (h *ImageHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if data.SourceUrl != "" && meta.SourceStatus == core.ProcessStatusInit {
		meta.SourceStatus = core.ProcessStatusUpdate
		meta.SourceError = ""
	}

	return nil
}

func (h *ImageHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == core.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == core.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *ImageHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	meta := node.Meta.(*ImageMeta)

	data := core.GetDownloadData()
	data.Filename = node.Name
	data.ContentType = meta.ContentType
	data.Stream = func(node *core.Node, w io.Writer) {
		_, err := h.Vault.Get(node.UniqueId(), w)
		core.PanicOnError(err)
	}

	return data
}

func (h *ImageHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *ImageHandler) StoreStream(node *core.Node, r io.Reader) (written int64, err error) {
	vaultmeta := core.GetVaultMetadata(node)

	meta := node.Meta.(*ImageMeta)
	meta.ContentType = "application/octet-stream"

	if written, err = h.Vault.Put(node.UniqueId(), vaultmeta, r); err != nil {
		core.PanicOnError(err)
	}

	return
}

type ImageDownloadListener struct {
	Vault      *vault.Vault
	HttpClient core.HttpClient
}

func (l *ImageDownloadListener) Handle(notification *pq.Notification, m core.NodeManager) (int, error) {
	reference, err := core.GetReferenceFromString(notification.Extra)

	if err != nil { // unable to parse the reference
		return core.PubSubListenContinue, nil
	}

	fmt.Printf("Download binary from uuid: %s\n", notification.Extra)
	node := m.Find(reference)

	if node == nil {
		fmt.Printf("Uuid does not exist: %s\n", notification.Extra)
		return core.PubSubListenContinue, nil
	}

	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == core.ProcessStatusDone {
		fmt.Printf("Nothing to update: %s\n", notification.Extra)

		return core.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(data.SourceUrl)

	if err != nil {
		meta.SourceStatus = core.ProcessStatusError
		meta.SourceError = "Unable to retrieve the remote file"
		m.Save(node, false)

		return core.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	vaultmeta := core.GetVaultMetadata(node)

	_, err = l.Vault.Put(node.UniqueId(), vaultmeta, resp.Body)

	if err != nil {
		return core.PubSubListenContinue, err
	}

	meta.ContentType = "application/octet-stream"
	meta.SourceStatus = core.ProcessStatusDone

	m.Save(node, false)

	return core.PubSubListenContinue, nil
}
