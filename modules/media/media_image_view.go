// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"errors"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	InvalidFitOptionsError = errors.New("Invalid fit options")
	InvalidWidthRangeError = errors.New("Invalid width range")
	WidthNotAllowedError   = errors.New("Width is not allowed")
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

func (m *MediaViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	values := request.HttpRequest.URL.Query()

	meta := node.Meta.(*ImageMeta)

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

	imageDst := resize.Resize(uint(width), 0, imageSrc, resize.Bicubic)

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

	var cropCenterX, cropCenterY int64
	var mode cutter.AnchorMode
	mode = cutter.Centered

	if len(options) == 4 {
		cropCenterX, err = strconv.ParseInt(options[2], 10, 0)
		if err != nil {
			return err
		}

		cropCenterY, err = strconv.ParseInt(options[3], 10, 0)
		if err != nil {
			return err
		}

		mode = cutter.TopLeft
	}

	imageSrc, err := m.getImage(node)

	if err != nil {
		return InvalidFitOptionsError
	}

	croppedImg, err := cutter.Crop(imageSrc, cutter.Config{
		Width:  int(width),
		Height: int(height),
		Anchor: image.Point{int(cropCenterX), int(cropCenterY)},
		Mode:   mode,
	})

	return m.encode(croppedImg, node, response.HttpResponse)
}

func (m *MediaViewHandler) getImage(node *base.Node) (image.Image, error) {
	source, _ := ioutil.TempFile(os.TempDir(), "gonode_image_resize_")
	defer source.Close()

	m.Vault.Get(node.UniqueId(), source)
	source.Seek(0, 0)

	meta := node.Meta.(*ImageMeta)

	if meta.ContentType == "image/jpeg" {
		return jpeg.Decode(source)
	}

	if meta.ContentType == "image/png" {
		return png.Decode(source)
	}

	return gif.Decode(source)
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
