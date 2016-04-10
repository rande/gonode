// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	sq "github.com/lann/squirrel"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/core/squirrel"
	"github.com/rande/gonode/modules/base"
)

const (
	OPERATION_OK = "OK"
	OPERATION_KO = "KO"
)

type ApiPager struct {
	Elements []interface{} `json:"elements"`
	Page     uint64        `json:"page"`
	PerPage  uint64        `json:"per_page"`
	Next     uint64        `json:"next"`
	Previous uint64        `json:"previous"`
}

type Api struct {
	Version    string
	Manager    base.NodeManager
	BaseUrl    string
	Serializer *base.Serializer
	Logger     *log.Logger
	Authorizer security.AuthorizationChecker
}

type ApiOperation struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a *Api) SelectBuilder(options *base.SelectOptions) sq.SelectBuilder {
	return a.Manager.SelectBuilder(options)
}

func (a *Api) Find(w io.Writer, query sq.SelectBuilder, page uint64, perPage uint64, options *base.AccessOptions) error {

	if options != nil && len(options.Roles) > 0 {
		value, _ := options.Roles.ToStringSlice()

		query = query.Where(squirrel.NewExprSlice(fmt.Sprintf("\"%s\" && ARRAY[" + sq.Placeholders(len(options.Roles)) + "]", "access"), value))
	}

	list := a.Manager.FindBy(query, (page - 1) * perPage, perPage + 1)

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
		a.Serializer.Serialize(b, e.Value.(*base.Node))

		message := json.RawMessage(b.Bytes())
		pager.Elements = append(pager.Elements, &message)

		counter++
	}

	base.Serialize(w, pager)

	return nil
}

func (a *Api) Save(r io.Reader, w io.Writer, options *base.AccessOptions) error {
	node := base.NewNode()

	err := a.Serializer.Deserialize(r, node)

	helper.PanicOnError(err)

	if a.Logger != nil {
		a.Logger.Printf("trying to save node.uuid=%s, node.type=%s", node.Uuid, node.Type)
	}

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		a.Logger.Printf("find uuid: %s", node.Uuid)

		helper.PanicUnless(node.Type == saved.Type, "Type mismatch")

		if options != nil {
			result, _ := a.Authorizer.IsGranted(options.Token, security.AttributesFromString(node.Access), node)

			if !result {
				return base.AccessForbiddenError
			}
		}

		if node.Deleted == true {
			return base.AlreadyDeletedError
		}

		if node.Revision != saved.Revision {
			return base.RevisionError
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
		base.Serialize(w, errors)

		return base.ValidationError
	}

	a.Manager.Save(node, true)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) Move(nodeUuid, parentUuid string, w io.Writer, options *base.AccessOptions) error {
	nodeReference, err := base.GetReferenceFromString(nodeUuid)

	if err != nil {
		return err
	}

	parentReference, err := base.GetReferenceFromString(parentUuid)

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

func (a *Api) FindOne(uuid string, w io.Writer, options *base.AccessOptions) error {
	reference, err := base.GetReferenceFromString(uuid)

	if err != nil {
		return base.NotFoundError
	}

	query := a.Manager.SelectBuilder(base.NewSelectOptions()).Where(sq.Eq{"uuid": reference.String()})

	return a.FindOneBy(query, w, options)
}

func (a *Api) FindOneBy(query sq.SelectBuilder, w io.Writer, options *base.AccessOptions) error {
	node := a.Manager.FindOneBy(query)

	if node == nil {
		return base.NotFoundError
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, security.AttributesFromString(node.Access), node)

		if !result {
			return base.AccessForbiddenError
		}
	}

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) RemoveOne(uuid string, w io.Writer, options *base.AccessOptions) error {
	reference, err := base.GetReferenceFromString(uuid)

	if err != nil {
		return err
	}

	node := a.Manager.Find(reference)

	if node == nil {
		return base.NotFoundError
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, security.AttributesFromString(node.Access), node)

		if !result {
			return base.AccessForbiddenError
		}
	}

	if node.Deleted {
		return base.AlreadyDeletedError
	}

	node, _ = a.Manager.RemoveOne(node)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) Remove(query sq.SelectBuilder, w io.Writer, options *base.AccessOptions) error {

	if options != nil && len(options.Roles) > 0 {
		value, _ := options.Roles.ToStringSlice()

		query = query.Where(squirrel.NewExprSlice(fmt.Sprintf("\"%s\" && ARRAY[" + sq.Placeholders(len(options.Roles)) + "]", "access"), value))
	}

	a.Manager.Remove(query)

	a.Find(w, query, 0, 0, options)

	return nil
}
