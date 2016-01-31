// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"github.com/rande/gonode/modules/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ImageHandler(t *testing.T) {
	a := assert.New(t)

	handler := &ImageHandler{}

	data, meta := handler.GetStruct()

	a.IsType(&ImageMeta{}, meta)
	a.IsType(&Image{}, data)

	a.Equal(meta.(*ImageMeta).SourceStatus, base.ProcessStatusInit)
	a.Equal(data.(*Image).SourceUrl, "")
}

func Test_ImageHandler_PreInsert(t *testing.T) {
	a := assert.New(t)

	node := base.NewNode()

	handler := &ImageHandler{}
	manager := &base.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusUpdate)

	// url, status done => no update
	node.Meta.(*ImageMeta).SourceStatus = base.ProcessStatusDone
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusDone)
}

func Test_ImageHandler_PreUpdate(t *testing.T) {
	a := assert.New(t)

	node := base.NewNode()

	handler := &ImageHandler{}
	manager := &base.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusUpdate)

	// url, status done => no update
	node.Meta.(*ImageMeta).SourceStatus = base.ProcessStatusDone
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusDone)
}

func Test_ImageHandler_PostUpdate(t *testing.T) {
	a := assert.New(t)

	node := base.NewNode()

	handler := &ImageHandler{}
	manager := &base.MockedManager{}
	manager.On("Notify", "media_file_download", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Meta.(*ImageMeta).SourceStatus = base.ProcessStatusUpdate

	handler.PostUpdate(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusUpdate)
}

func Test_ImageHandler_PostInsert(t *testing.T) {
	a := assert.New(t)

	node := base.NewNode()

	handler := &ImageHandler{}
	manager := &base.MockedManager{}
	manager.On("Notify", "media_file_download", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Meta.(*ImageMeta).SourceStatus = base.ProcessStatusUpdate

	handler.PostInsert(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Meta.(*ImageMeta).SourceStatus, base.ProcessStatusUpdate)
}
