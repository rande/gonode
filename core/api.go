// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	sq "github.com/lann/squirrel"
	"io"
	"log"
)

type SearchForm struct {
	Page     int64               `schema:"page"`
	PerPage  int64               `schema:"per_page"`
	OrderBy  []string            `schema:"order_by"`
	Uuid     string              `schema:"uuid"`
	Type     []string            `schema:"type"`
	Name     string              `schema:"name"`
	Slug     string              `schema:"slug"`
	Data     map[string][]string `schema:"data"`
	Meta     map[string][]string `schema:"meta"`
	Status   []string            `schema:"status"`
	Weight   []string            `schema:"weight"`
	Revision string              `schema:"revision"`
	//	CreatedAt  time.Time          `schema:"created_at"`
	//	UpdatedAt  time.Time          `schema:"updated_at"`
	Enabled string `schema:"enabled"`
	Deleted string `schema:"deleted"`
	Current string `schema:"current"`
	//	Parents    []Reference        `schema:"parents"`
	UpdatedBy  []string `schema:"updated_by"`
	CreatedBy  []string `schema:"created_by"`
	ParentUuid []string `schema:"parent_uuid"`
	SetUuid    []string `schema:"set_uuid"`
	Source     []string `schema:"source"`
}

func GetSearchForm() *SearchForm {
	return &SearchForm{
		Data:    make(map[string][]string),
		Meta:    make(map[string][]string),
		OrderBy: []string{"updated_at,ASC"},
	}
}

type ApiPager struct {
	Elements []interface{} `json:"elements"`
	Page     uint64        `json:"page"`
	PerPage  uint64        `json:"per_page"`
	Next     uint64        `json:"next"`
	Previous uint64        `json:"previous"`
}

type Api struct {
	Version    string
	Manager    NodeManager
	BaseUrl    string
	Serializer *Serializer
	Logger     *log.Logger
}

func (a *Api) SelectBuilder() sq.SelectBuilder {
	return a.Manager.SelectBuilder()
}

func (a *Api) Find(w io.Writer, query sq.SelectBuilder, page uint64, perPage uint64) error {
	list := a.Manager.FindBy(query, (page-1)*perPage, perPage+1)

	pager := &ApiPager{
		Page:    page,
		PerPage: perPage,
	}

	pager.Elements = make([]interface{}, 0)

	if page > 1 {
		pager.Previous = page - 1
	}

	for e := list.Front(); e != nil; e = e.Next() {
		b := bytes.NewBuffer([]byte{})
		a.Serializer.Serialize(b, e.Value.(*Node))

		message := json.RawMessage(b.Bytes())
		pager.Elements = append(pager.Elements, &message)
	}

	Serialize(w, pager)

	return nil
}

func (a *Api) Save(r io.Reader, w io.Writer) error {
	node := NewNode()

	err := a.Serializer.Deserialize(r, node)

	PanicOnError(err)

	if a.Logger != nil {
		a.Logger.Printf("trying to save node.uuid=%s, node.type=%s", node.Uuid, node.Type)
	}

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		a.Logger.Printf("find uuid: %s", node.Uuid)

		PanicUnless(node.Type == saved.Type, "Type mismatch")
		PanicIf(saved.Deleted, "Cannot save a deleted node, restore it first ...")

		if node.Revision != saved.Revision {
			return RevisionError
		}

		node.id = saved.id
	} else if a.Logger != nil {
		a.Logger.Printf("cannot find uuid: %s, create a new one", node.Uuid)
	}

	if a.Logger != nil {
		a.Logger.Printf("saving node.id=%d, node.uuid=%s", node.id, node.Uuid)
	}

	if ok, errors := a.Manager.Validate(node); !ok {
		Serialize(w, errors)

		return ValidationError
	}

	a.Manager.Save(node)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) FindOne(uuid string, w io.Writer) error {

	node := a.Manager.Find(GetReferenceFromString(uuid))

	if node == nil {
		return NotFoundError
	}

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) RemoveOne(uuid string, w io.Writer) error {
	node := a.Manager.Find(GetReferenceFromString(uuid))

	if node == nil {
		return NotFoundError
	}

	if node.Deleted {
		return AlreadyDeletedError
	}

	node, _ = a.Manager.RemoveOne(node)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) Remove(b sq.SelectBuilder, w io.Writer) error {
	a.Manager.Remove(b)

	a.Find(w, b, 0, 0)

	return nil
}
