// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
)

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
	Vault  *vault.Vault
	Logger *log.Logger
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

func (h *ImageHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return HandleImageReader(node, h.Vault, r, h.Logger)
}

type ImageDownloadListener struct {
	Vault      *vault.Vault
	HttpClient base.HttpClient
	Logger     *log.Logger
}

func (l *ImageDownloadListener) Handle(notification *pq.Notification, m base.NodeManager) (int, error) {
	reference, err := base.GetReferenceFromString(notification.Extra)

	if err != nil {
		l.Logger.WithFields(log.Fields{
			"module":    "media.downloader",
			"node_uuid": notification.Extra,
		}).Debug("Unable to parse reference")

		// unable to parse the reference
		return base.PubSubListenContinue, nil
	}

	if l.Logger != nil {
		l.Logger.WithFields(log.Fields{
			"module":    "media.downloader",
			"node_uuid": notification.Extra,
		}).Debug("Download binary from uuid")
	}

	node := m.Find(reference)

	if node == nil {
		if l.Logger != nil {
			l.Logger.WithFields(log.Fields{
				"module":    "media.downloader",
				"node_uuid": notification.Extra,
			}).Info("Unable to download file, uuid does not exist")
		}

		return base.PubSubListenContinue, nil
	}

	data := node.Data.(*Image)
	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus == base.ProcessStatusDone {
		if l.Logger != nil {
			l.Logger.WithFields(log.Fields{
				"module":             "media.downloader",
				"node_uuid":          notification.Extra,
				"meta_source_status": base.ProcessStatusDone,
			}).Warn("Stop downloading process, already done! (race condition ?)")
		}

		return base.PubSubListenContinue, nil
	}

	resp, err := l.HttpClient.Get(data.SourceUrl)

	if err != nil {
		meta.SourceStatus = base.ProcessStatusError
		meta.SourceError = "Unable to retrieve the remote file"

		if l.Logger != nil {
			l.Logger.WithFields(log.Fields{
				"module":             "media.downloader",
				"node_uuid":          notification.Extra,
				"error":              err.Error(),
				"meta_source_status": base.ProcessStatusError,
			}).Warn("Unable to retrieve the remote file")
		}

		m.Save(node, false)

		return base.PubSubListenContinue, err
	}

	defer resp.Body.Close()

	if _, err = HandleImageReader(node, l.Vault, resp.Body, l.Logger); err != nil {
		meta.SourceStatus = base.ProcessStatusError
		meta.SourceError = "Unable to analyse the image"

		m.Save(node, false)

		return base.PubSubListenContinue, err
	}

	meta.SourceStatus = base.ProcessStatusDone

	m.Save(node, false)

	if l.Logger != nil {
		l.Logger.WithFields(log.Fields{
			"module":    "media.downloader",
			"node_uuid": notification.Extra,
		}).Debug("File downloaded!")
	}

	return base.PubSubListenContinue, nil
}
