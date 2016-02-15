// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/goexif/exif"
	"github.com/rande/goexif/mknote"
	"github.com/rande/goexif/tiff"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ExifMeta map[string]string

type ExifWalker struct {
	Meta ExifMeta
}

func (w *ExifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	w.Meta[string(name)] = strings.Trim(tag.String(), "\"")

	return nil
}

func GetExif(r io.Reader) (ExifMeta, error) {
	// try to find exif information
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(r)

	if err != nil {
		return nil, err
	}

	e := &ExifWalker{
		Meta: make(map[string]string, 0),
	}

	if err != nil {
		return nil, err
	}

	x.Walk(e)

	return e.Meta, nil
}

func HandleImageReader(node *base.Node, v *vault.Vault, r io.Reader, logger *log.Logger) (int64, error) {

	vaultmeta := base.GetVaultMetadata(node)
	meta := node.Meta.(*ImageMeta)

	f, _ := ioutil.TempFile(os.TempDir(), "gonode_media_downloader_")
	path := f.Name()

	if logger != nil {
		logger.WithFields(log.Fields{
			"module":    "media.handle_image_reader",
			"node_uuid": node.Uuid.String(),
			"path":      path,
		}).Debug("Start handling io.Reader to store image")
	}

	defer func() {
		err := f.Close()

		if err != nil {
			if logger != nil {
				logger.WithFields(log.Fields{
					"module":    "media.downloader",
					"node_uuid": node.Uuid.String(),
					"error":     err.Error(),
					"path":      path,
				}).Warn("Unable to close temporary file")
			}
		}

		err = os.Remove(path)

		if err != nil {
			if logger != nil {
				logger.WithFields(log.Fields{
					"module":    "media.downloader",
					"node_uuid": node.Uuid.String(),
					"error":     err.Error(),
					"path":      path,
				}).Warn("Unable to delete temporary file")
			}
		}
	}()

	written, err := io.Copy(f, r)
	if err != nil {
		if logger != nil {
			logger.WithFields(log.Fields{
				"module":    "media.downloader",
				"node_uuid": node.Uuid.String(),
				"error":     err.Error(),
				"path":      path,
			}).Warn("Unable to copy io.Reader's buffer to temporary file")
		}

		return written, err
	}

	f.Seek(0, 0)

	_, err = v.Put(node.UniqueId(), vaultmeta, f)

	if err != nil {
		if logger != nil {
			logger.WithFields(log.Fields{
				"module":    "media.downloader",
				"node_uuid": node.Uuid.String(),
				"error":     err.Error(),
				"path":      path,
			}).Warn("Unable to put the file into the vault")
		}

		return written, err
	}

	f.Seek(0, 0)

	d := make([]byte, 500)
	f.Read(d)

	meta.ContentType = http.DetectContentType(d)

	if meta.ContentType == "image/jpeg" {
		f.Seek(0, 0)
		em, err := GetExif(f)

		if err != nil {
			if err != nil {
				logger.WithFields(log.Fields{
					"module":    "media.downloader",
					"node_uuid": node.Uuid.String(),
					"error":     err.Error(),
					"path":      path,
				}).Info("Unable to decode exif data")
			}

			return written, err
		}

		meta.Exif = em
	}

	if logger != nil {
		logger.WithFields(log.Fields{
			"module":    "media.handle_image_reader",
			"node_uuid": node.Uuid.String(),
			"path":      path,
		}).Debug("End handling io.Reader to store image")
	}

	return written, nil
}
