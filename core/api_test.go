package core

import (
	"bytes"
	"container/list"
	"encoding/json"
	sq "github.com/lann/squirrel"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func Test_ApiPager_Serialization(t *testing.T) {

	sb := sq.Select("id, name").From("test_nodes").PlaceholderFormat(sq.Dollar)

	list := list.New()
	node1 := NewNode()
	node1.Type = "image"
	node1.CreatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")
	node1.UpdatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")

	list.PushBack(node1)

	node2 := NewNode()
	node2.Type = "video"
	node2.CreatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")
	node2.UpdatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")

	list.PushBack(node2)

	manager := &MockedManager{}
	manager.On("SelectBuilder").Return(sb)
	manager.On("FindBy", sb, uint64(0), uint64(11)).Return(list)

	api := &Api{
		Version:    "1",
		Manager:    manager,
		Serializer: NewSerializer(),
	}

	b := bytes.NewBuffer([]byte{})

	assert.Equal(t, sb, api.SelectBuilder())

	api.Find(b, api.SelectBuilder(), uint64(1), uint64(10))

	var out bytes.Buffer

	json.Indent(&out, b.Bytes(), "", "    ")

	data, err := ioutil.ReadFile("../test/fixtures/pager_results.json")

	assert.Equal(t, string(data[:]), out.String())
}
