// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"io"
	"net/http"
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

func (h *ImageHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Image{}, &ImageMeta{
		SourceStatus: base.ProcessStatusInit,
	}
}

func (h *ImageHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if data.SourceUrl != "" && meta.SourceStatus == base.ProcessStatusInit {
		meta.SourceStatus = base.ProcessStatusUpdate
		meta.SourceError = ""
	}

	return nil
}

func (h *ImageHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if data.SourceUrl != "" && meta.SourceStatus == base.ProcessStatusInit {
		meta.SourceStatus = base.ProcessStatusUpdate
		meta.SourceError = ""
	}

	return nil
}

func (h *ImageHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == base.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == base.ProcessStatusUpdate {
		m.Notify("media_file_download", node.Uuid.String())
	}

	return nil
}

func (h *ImageHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {

}

func (h *ImageHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	meta := node.Meta.(*ImageMeta)

	data := base.GetDownloadData()
	data.Filename = node.Name
	data.ContentType = meta.ContentType
	data.Stream = func(node *base.Node, w io.Writer) {
		_, err := h.Vault.Get(node.UniqueId(), w)
		helper.PanicOnError(err)
	}

	return data
}

func (h *ImageHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *ImageHandler) StoreStream(node *base.Node, r io.Reader) (written int64, err error) {
	vaultmeta := base.GetVaultMetadata(node)

	meta := node.Meta.(*ImageMeta)
	meta.ContentType = "application/octet-stream"

	if written, err = h.Vault.Put(node.UniqueId(), vaultmeta, r); err != nil {
		helper.PanicOnError(err)
	}

	return
}

type ImageDownloadListener struct {
	Vault      *vault.Vault
	HttpClient base.HttpClient
}

func (l *ImageDownloadListener) Handle(notification *pq.Notification, m base.NodeManager) (int, error) {
	reference, err := base.GetReferenceFromString(notification.Extra)

	if err != nil { // unable to parse the reference
		return base.PubSubListenContinue, nil
	}

	fmt.Printf("Download binary from uuid: %s\n", notification.Extra)
	node := m.Find(reference)

	if node == nil {
		fmt.Printf("Uuid does not exist: %s\n", notification.Extra)
		return base.PubSubListenContinue, nil
	}

	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == base.ProcessStatusDone {
		fmt.Printf("Nothing to update: %s\n", notification.Extra)

		return base.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(data.SourceUrl)

	if err != nil {
		meta.SourceStatus = base.ProcessStatusError
		meta.SourceError = "Unable to retrieve the remote file"
		m.Save(node, false)

		return base.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	vaultmeta := base.GetVaultMetadata(node)

	r := &helper.PartialReader{
		Reader: resp.Body,
		Size:   500,
	}

	_, err = l.Vault.Put(node.UniqueId(), vaultmeta, r)

	if err != nil {
		return base.PubSubListenContinue, err
	}

	meta.ContentType = http.DetectContentType(r.Data)
	meta.SourceStatus = base.ProcessStatusDone

	m.Save(node, false)

	return base.PubSubListenContinue, nil
}
