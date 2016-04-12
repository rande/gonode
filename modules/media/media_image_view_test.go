// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"github.com/stretchr/testify/assert"
)

func Test_ContainsSize(t *testing.T) {
	assert.True(t, ContainsSize(200, []uint{}))
	assert.True(t, ContainsSize(200, []uint{100, 200, 400}))

	assert.False(t, ContainsSize(200, []uint{100, 300}))
}

func GetDriver() vault.VaultDriver {
	fv, err := os.Open("../../test/fixtures/photo.jpg.vault")

	helper.PanicOnError(err)

	driver := &vault.MockedDriver{}
	driver.
		// default string for 111...111 uuid
		On("GetReader", "eb/60/6938046e0b96477b491c35ab8ce174ce96cfef588c827508a14822e16939.bin.vault").
		Return(fv, nil)

	f, err := os.Open("../../test/fixtures/photo.jpg")
	helper.PanicOnError(err)
	driver.
		// default string for 111...111 uuid
		On("GetReader", "eb/60/6938046e0b96477b491c35ab8ce174ce96cfef588c827508a14822e16939.bin").
		Return(f, nil)

	return driver
}

func Test_MediaViewHandler_NoError_And_NoParam(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", res.Header().Get("Content-Type"))
	assert.True(t, len(res.Body.Bytes()) > 1024)
}

func Test_MediaViewHandler_VaultError(t *testing.T) {
	driver := &vault.MockedDriver{}
	driver.
		// default string for 111...111 uuid
		On("GetReader", "eb/60/6938046e0b96477b491c35ab8ce174ce96cfef588c827508a14822e16939.bin.vault").
		Return(nil, errors.New("error"))

	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: driver,
		},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.Error(t, err)
	assert.Equal(t, "image/jpeg", res.Header().Get("Content-Type"))
}

func Test_MediaViewHandler_InvalidResize_Width(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth: 300,
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mr=3000", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.Error(t, err)
	assert.Equal(t, err, InvalidWidthRangeError)
}

func Test_MediaViewHandler_Invalid_Crop_Width(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth: 300,
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mf=3000,3000", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.Error(t, err)
	assert.Equal(t, err, InvalidWidthRangeError)
}

func Test_MediaViewHandler_Invalid_Resize_NotAllowedWidth(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth:      300,
		AllowedWidths: []uint{100, 150},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mr=120", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.Error(t, err)
	assert.Equal(t, err, WidthNotAllowedError)
}

func Test_MediaViewHandler_Invalid_Crop_NotAllowedWidth(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth:      300,
		AllowedWidths: []uint{100, 150},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mf=120,120", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.Error(t, err)
	assert.Equal(t, err, WidthNotAllowedError)
}

func Test_MediaViewHandler_Resize(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth:      300,
		AllowedWidths: []uint{100, 150},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mr=150", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", res.Header().Get("Content-Type"))
	assert.True(t, len(res.Body.Bytes()) > 1024)
}

func Test_MediaViewHandler_Crop(t *testing.T) {
	n := base.NewNode()
	n.Meta = &ImageMeta{
		ContentType:  "image/jpeg",
		SourceStatus: base.ProcessStatusDone,
	}

	m := &MediaViewHandler{
		Vault: &vault.Vault{
			Algo:   "no_op",
			Driver: GetDriver(),
		},
		MaxWidth:      300,
		AllowedWidths: []uint{100, 150},
	}

	req, _ := http.NewRequest("GET", "/we-don-t-care?mf=150,150", nil)

	request := &base.ViewRequest{
		HttpRequest: req,
	}

	res := httptest.NewRecorder()
	response := &base.ViewResponse{
		HttpResponse: res,
	}

	err := m.Execute(n, request, response)

	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", res.Header().Get("Content-Type"))
	assert.True(t, len(res.Body.Bytes()) > 1024)
}
