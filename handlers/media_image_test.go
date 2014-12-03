package handlers

import (
	nm "github.com/rande/gonode/test/mock"
	nc "github.com/rande/gonode/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ImageHandler(t *testing.T) {
	a := assert.New(t)

	handler := &ImageHandler{}

	data, meta := handler.GetStruct()

	a.IsType(&ImageMeta{}, meta)
	a.IsType(&Image{}, data)

	a.Equal(data.(*Image).SourceStatus, nc.ProcessStatusInit)
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
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreInsert(node, manager)
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusUpdate)

	// url, status done => no update
	node.Data.(*Image).SourceStatus = nc.ProcessStatusDone
	handler.PreInsert(node, manager)
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusDone)
}

func Test_ImageHandler_PreUpdate(t *testing.T) {
	a := assert.New(t)

	node := nc.NewNode()

	handler := &ImageHandler{}
	manager := &nc.PgNodeManager{}

	node.Data, node.Meta = handler.GetStruct()

	// no url => keep status
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusInit)

	// url => update status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusUpdate)

	// url, status done => no update
	node.Data.(*Image).SourceStatus = nc.ProcessStatusDone
	handler.PreUpdate(node, manager)
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusDone)
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
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Data.(*Image).SourceStatus = nc.ProcessStatusUpdate

	handler.PostUpdate(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusUpdate)
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
	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusInit)

	// url => keep status
	node.Data.(*Image).SourceUrl = "http://myimage.com/salut.jpg"
	node.Data.(*Image).SourceStatus = nc.ProcessStatusUpdate

	handler.PostInsert(node, manager)

	manager.AssertCalled(t, "Notify", "media_file_download", node.Uuid.String())

	a.Equal(node.Data.(*Image).SourceStatus, nc.ProcessStatusUpdate)
}
