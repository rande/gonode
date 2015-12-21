// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	sq "github.com/lann/squirrel"
	"github.com/rande/gonode/core"
	"io"
	"log"
)

const (
	OPERATION_OK = "OK"
	OPERATION_KO = "KO"
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
	Manager    core.NodeManager
	BaseUrl    string
	Serializer *core.Serializer
	Logger     *log.Logger
}

type ApiOperation struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a *Api) SelectBuilder(options *core.SelectOptions) sq.SelectBuilder {
	return a.Manager.SelectBuilder(options)
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

	counter := uint64(0)
	for e := list.Front(); e != nil; e = e.Next() {
		if counter == perPage {
			pager.Next = page + 1
			break
		}

		b := bytes.NewBuffer([]byte{})
		a.Serializer.Serialize(b, e.Value.(*core.Node))

		message := json.RawMessage(b.Bytes())
		pager.Elements = append(pager.Elements, &message)

		counter++
	}

	core.Serialize(w, pager)

	return nil
}

func (a *Api) Save(r io.Reader, w io.Writer) error {
	node := core.NewNode()

	err := a.Serializer.Deserialize(r, node)

	core.PanicOnError(err)

	if a.Logger != nil {
		a.Logger.Printf("trying to save node.uuid=%s, node.type=%s", node.Uuid, node.Type)
	}

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		a.Logger.Printf("find uuid: %s", node.Uuid)

		core.PanicUnless(node.Type == saved.Type, "Type mismatch")
		core.PanicIf(saved.Deleted, "Cannot save a deleted node, restore it first ...")

		if node.Revision != saved.Revision {
			return core.RevisionError
		}

		node.Id = saved.Id

		// we cannot overwrite the Parents, Or the ParentUuid, need to use the http API
		node.Parents = saved.Parents
		node.ParentUuid = saved.ParentUuid

	} else if a.Logger != nil {
		a.Logger.Printf("cannot find uuid: %s, create a new one", node.Uuid)
	}

	if a.Logger != nil {
		a.Logger.Printf("saving node.id=%d, node.uuid=%s", node.Id, node.Uuid)
	}

	if ok, errors := a.Manager.Validate(node); !ok {
		core.Serialize(w, errors)

		return core.ValidationError
	}

	a.Manager.Save(node, true)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) Move(nodeUuid, parentUuid string, w io.Writer) error {

	nodeReference, err := core.GetReferenceFromString(nodeUuid)

	if err != nil {
		return err
	}

	parentReference, err := core.GetReferenceFromString(parentUuid)

	if err != nil {
		return err
	}

	affectedNodes, err := a.Manager.Move(nodeReference, parentReference)

	if err != nil {
		return err
	}

	a.Serializer.Serialize(w, &ApiOperation{
		Status:  OPERATION_OK,
		Message: fmt.Sprintf("Node altered: %d", affectedNodes),
	})

	return nil
}

func (a *Api) FindOne(uuid string, w io.Writer) error {
	reference, err := core.GetReferenceFromString(uuid)

	if err != nil {
		return core.NotFoundError
	}

	node := a.Manager.Find(reference)

	if node == nil {
		return core.NotFoundError
	}

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) FindOneBy(query sq.SelectBuilder, w io.Writer) error {

	node := a.Manager.FindOneBy(query)

	if node == nil {
		return core.NotFoundError
	}

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) RemoveOne(uuid string, w io.Writer) error {
	reference, err := core.GetReferenceFromString(uuid)

	if err != nil {
		return err
	}

	node := a.Manager.Find(reference)

	if node == nil {
		return core.NotFoundError
	}

	if node.Deleted {
		return core.AlreadyDeletedError
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
