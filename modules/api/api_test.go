// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"github.com/gorilla/schema"
	sq "github.com/lann/squirrel"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func Test_ApiPager_Serialization(t *testing.T) {
	sb := sq.Select("id, name").From("test_nodes").PlaceholderFormat(sq.Dollar)

	list := list.New()
	node1 := base.NewNode()
	node1.Type = "image"
	node1.CreatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")
	node1.UpdatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")

	list.PushBack(node1)

	node2 := base.NewNode()
	node2.Type = "video"
	node2.CreatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")
	node2.UpdatedAt, _ = time.Parse(time.RFC3339Nano, "2015-06-15T10:23:08.698707603+02:00")

	list.PushBack(node2)

	options := base.NewSelectOptions()
	manager := &base.MockedManager{}
	manager.On("SelectBuilder", options).Return(sb)
	manager.On("FindBy", sb, uint64(0), uint64(11)).Return(list)

	api := &Api{
		Version:    "1",
		Manager:    manager,
		Serializer: base.NewSerializer(),
	}

	b := bytes.NewBuffer([]byte{})

	assert.Equal(t, sb, api.SelectBuilder(options))

	api.Find(b, api.SelectBuilder(options), uint64(1), uint64(10))

	var out bytes.Buffer

	json.Indent(&out, b.Bytes(), "", "    ")

	data, err := ioutil.ReadFile("../../test/fixtures/pager_results.json")

	helper.PanicOnError(err)

	assert.Equal(t, string(data[:]), out.String())
}

func Test_ApiPager_Deserialization(t *testing.T) {
	data, _ := ioutil.ReadFile("../../test/fixtures/pager_results.json")

	p := &ApiPager{}

	json.Unmarshal(data, p)

	assert.Equal(t, uint64(10), p.PerPage)
	assert.Equal(t, uint64(1), p.Page)
	assert.Equal(t, 2, len(p.Elements))
	assert.Equal(t, uint64(0), p.Next)
	assert.Equal(t, uint64(0), p.Previous)

	// the Element is a [string]interface so we need to convert it back to []byte
	// and then unmarshal again with the correct structure
	raw, _ := json.Marshal(p.Elements[0])

	n := base.NewNode()
	json.Unmarshal(raw, n)

	assert.Equal(t, "image", n.Type)
	assert.False(t, n.Deleted)
}

func Test_OrderBy_Form(t *testing.T) {
	form := search.GetHttpSearchForm()

	values := map[string][]string{
		"order_by": {
			"updated_at,ASC",
			"name,DESC",
		},
	}

	decoder := schema.NewDecoder()
	decoder.Decode(form, values)

	assert.Equal(t, form.OrderBy, []string{"updated_at,ASC", "name,DESC"})
}
