// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
)

var (
	InvalidFitOptionsError = errors.New("Invalid fit options")
	InvalidWidthRangeError = errors.New("Invalid width range")
	WidthNotAllowedError   = errors.New("Width is not allowed")
	InvalidProcessStatus   = errors.New("Invalid process status")
)

func ContainsSize(size uint, allowed []uint) bool {
	if len(allowed) == 0 {
		return true
	}

	for _, v := range allowed {
		if v == size {
			return true
		}
	}

	return false
}

type MediaViewHandler struct {
	Vault         *vault.Vault
	AllowedWidths []uint
	MaxWidth      uint
}

func (m *MediaViewHandler) Support(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) bool {
	return request.Format == "jpg" || request.Format == "gif" || request.Format == "png"
}

func (m *MediaViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	values := request.HttpRequest.URL.Query()

	meta := node.Meta.(*ImageMeta)

	if meta.SourceStatus != base.ProcessStatusInit && meta.SourceStatus != base.ProcessStatusDone {
		return InvalidProcessStatus
	}

	response.HttpResponse.Header().Set("Content-Type", meta.ContentType)

	for _, format := range []string{"image/jpeg", "image/gif", "image/png"} {
		if format == meta.ContentType { // only resize supported format
			if _, ok := values["mr"]; ok { // ask for binary content
				return m.imageResize(node, request, response)
			} else if _, ok := values["mf"]; ok {
				return m.imageFit(node, request, response)
			}
		}
	}

	_, err := m.Vault.Get(node.UniqueId(), response.HttpResponse)

	return err
}

func (m *MediaViewHandler) imageResize(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	width, err := strconv.ParseUint(request.HttpRequest.URL.Query().Get("mr"), 10, 0)

	if err != nil {
		return err
	}

	if m.MaxWidth > 0 && uint(width) > m.MaxWidth {
		return InvalidWidthRangeError
	}

	if !ContainsSize(uint(width), m.AllowedWidths) {
		return WidthNotAllowedError
	}

	imageSrc, err := m.getImage(node)

	if err != nil {
		return err
	}

	imageDst := imaging.Resize(imageSrc, int(width), 0, imaging.Lanczos)

	return m.encode(imageDst, node, response.HttpResponse)
}

func (m *MediaViewHandler) imageFit(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	options := strings.Split(request.HttpRequest.URL.Query().Get("mf"), ",")

	if len(options) != 2 && len(options) != 4 {
		return InvalidFitOptionsError
	}

	width, err := strconv.ParseInt(options[0], 10, 0)
	if err != nil {
		return err
	}

	height, err := strconv.ParseInt(options[1], 10, 0)
	if err != nil {
		return err
	}

	if m.MaxWidth > 0 && uint(width) > m.MaxWidth {
		return InvalidWidthRangeError
	}

	if !ContainsSize(uint(width), m.AllowedWidths) {
		return WidthNotAllowedError
	}

	imageSrc, err := m.getImage(node)

	if err != nil {
		return InvalidFitOptionsError
	}

	croppedImg := imaging.Fill(imageSrc, int(width), int(height), imaging.Center, imaging.Lanczos)

	return m.encode(croppedImg, node, response.HttpResponse)
}

func (m *MediaViewHandler) getImage(node *base.Node) (image.Image, error) {
	source, _ := ioutil.TempFile(os.TempDir(), "gonode_image_resize_")
	defer source.Close()

	m.Vault.Get(node.UniqueId(), source)
	source.Seek(0, 0)

	meta := node.Meta.(*ImageMeta)

	var img image.Image
	var err error
	if meta.ContentType == "image/jpeg" {
		img, err = jpeg.Decode(source)
	} else if meta.ContentType == "image/png" {
		img, err = png.Decode(source)
	} else {
		img, err = gif.Decode(source)
	}

	if _, ok := meta.Exif["Orientation"]; ok {
		// http://piexif.readthedocs.org/en/latest/sample.html#rotate-image-by-exif-orientation
		switch meta.Exif["Orientation"] {
		case "2":
			img = imaging.FlipH(img)
		case "3":
			img = imaging.Rotate180(img)
		case "4":
			img = imaging.FlipH(imaging.Rotate180(img))
		case "5":
			img = imaging.FlipH(imaging.Rotate90(img))
		case "6":
			img = imaging.Rotate270(img)
		case "7":
			img = imaging.FlipH(imaging.Rotate90(img))
		case "8":
			img = imaging.Rotate90(img)
		}
	}

	return img, err
}

func (m *MediaViewHandler) encode(image image.Image, node *base.Node, w io.Writer) error {
	meta := node.Meta.(*ImageMeta)

	if meta.ContentType == "image/jpeg" {
		return jpeg.Encode(w, image, nil)
	}

	if meta.ContentType == "image/png" {
		return png.Encode(w, image)
	}

	return gif.Encode(w, image, nil)
}
