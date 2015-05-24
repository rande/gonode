package handlers

import (
	nc "github.com/rande/gonode/core"
	nm "github.com/rande/gonode/test/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ImageHandler(t *testing.T) {
	a := assert.New(t)

	handler := &ImageHandler{}

	data, meta := handler.GetStruct()

	a.IsType(&ImageMeta{}, meta)
	a.IsType(&Image{}, data)

	a.Equal(meta.(*ImageMeta).SourceStatus, nc.ProcessStatusInit)
	a.Equal(data.(*Image).SourceUrl, "")
}

func Test_ImageHandler_PreInsert(t *testing.T) {
	a := assert.New(t)

	node := nc.NewNode()

	handler := &ImageHandler{}
	manager := &nc.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusUpdate)

	// url, status done => no update
	node.Meta.(*ImageMeta).SourceStatus = nc.ProcessStatusDone
	handler.PreInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusDone)
}

func Test_ImageHandler_PreUpdate(t *testing.T) {
	a := assert.New(t)

	node := nc.NewNode()

	handler := &ImageHandler{}
	manager := &nc.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusUpdate)

	// url, status done => no update
	node.Meta.(*ImageMeta).SourceStatus = nc.ProcessStatusDone
	handler.PreUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusDone)
}

func Test_ImageHandler_PostUpdate(t *testing.T) {
	a := assert.New(t)

	node := nc.NewNode()

	handler := &ImageHandler{}
	manager := &nm.MockedManager{}
	manager.On("Notify", "media_file_download", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostUpdate(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Meta.(*ImageMeta).SourceStatus = nc.ProcessStatusUpdate

	handler.PostUpdate(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusUpdate)
}

func Test_ImageHandler_PostInsert(t *testing.T) {
	a := assert.New(t)

	node := nc.NewNode()

	handler := &ImageHandler{}
	manager := &nm.MockedManager{}
	manager.On("Notify", "media_file_download", node.Uuid.String()).Return()

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PostInsert(node, manager)
	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Meta.(*ImageMeta).SourceStatus = nc.ProcessStatusUpdate

	handler.PostInsert(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Meta.(*ImageMeta).SourceStatus, nc.ProcessStatusUpdate)
}
