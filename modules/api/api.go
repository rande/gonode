// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

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

func (a *Api) Find(query sq.SelectBuilder, page uint64, perPage uint64, options *base.AccessOptions) (*ApiPager, error) {
	if options != nil && len(options.Roles) > 0 {
		value, _ := options.Roles.ToStringSlice()

		query = query.Where(squirrel.NewExprSlice(fmt.Sprintf("\"%s\" && ARRAY["+sq.Placeholders(len(options.Roles))+"]", "access"), value))
	}

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

		pager.Elements = append(pager.Elements, e.Value.(*base.Node))

		counter++
	}

	return pager, nil
}

func (a *Api) Save(node *base.Node, options *base.AccessOptions) (*base.Node, base.Errors, error) {
	if a.Logger != nil {
		a.Logger.Printf("trying to save node.uuid=%s, node.type=%s", node.Uuid, node.Type)
	}

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		a.Logger.Printf("find uuid: %s", node.Uuid)

		helper.PanicUnless(node.Type == saved.Type, "Type mismatch")

		if options != nil {
			result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

			if !result {
				return nil, nil, base.AccessForbiddenError
			}
		}

		if node.Deleted == true {
			return nil, nil, base.AlreadyDeletedError
		}

		if node.Revision != saved.Revision {
			return nil, nil, base.RevisionError
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
		return nil, errors, base.ValidationError
	}

	node, err := a.Manager.Save(node, true)

	return node, nil, err
}

func (a *Api) Move(nodeUuid, parentUuid string, options *base.AccessOptions) (*ApiOperation, error) {
	// handle node
	nodeReference, err := base.GetReferenceFromString(nodeUuid)

	if err != nil {
		return nil, err
	}

	if node := a.Manager.Find(nodeReference); node == nil {
		return nil, base.NotFoundError
	} else if result, _ := a.Authorizer.IsGranted(options.Token, nil, node); !result {
		return nil, base.AccessForbiddenError
	}

	// parent node
	parentReference, err := base.GetReferenceFromString(parentUuid)

	if err != nil {
		return nil, err
	}

	if parent := a.Manager.Find(parentReference); parent == nil {
		return nil, base.NotFoundError
	} else if result, _ := a.Authorizer.IsGranted(options.Token, nil, parent); !result {
		return nil, base.AccessForbiddenError
	}

	// move node
	if affectedNodes, err := a.Manager.Move(nodeReference, parentReference); err != nil {
		return nil, err
	} else {
		return &ApiOperation{
			Status:  OPERATION_OK,
			Message: fmt.Sprintf("Node altered: %d", affectedNodes),
		}, nil
	}
}

func (a *Api) FindOne(uuid string, options *base.AccessOptions) (*base.Node, error) {
	reference, err := base.GetReferenceFromString(uuid)

	if err != nil {
		return nil, base.NotFoundError
	}

	query := a.Manager.SelectBuilder(base.NewSelectOptions()).Where(sq.Eq{"uuid": reference.String()})

	return a.FindOneBy(query, options)
}

func (a *Api) FindOneBy(query sq.SelectBuilder, options *base.AccessOptions) (*base.Node, error) {
	node := a.Manager.FindOneBy(query)

	if node == nil {
		return nil, base.NotFoundError
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

		if !result {
			return nil, base.AccessForbiddenError
		}
	}

	return node, nil
}

func (a *Api) RemoveOne(uuid string, options *base.AccessOptions) (*base.Node, error) {
	reference, err := base.GetReferenceFromString(uuid)

	if err != nil {
		return nil, err
	}

	node := a.Manager.Find(reference)

	if node == nil {
		return nil, base.NotFoundError
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

		if !result {
			return nil, base.AccessForbiddenError
		}
	}

	if node.Deleted {
		return nil, base.AlreadyDeletedError
	}

	return a.Manager.RemoveOne(node)
}

func (a *Api) Remove(query sq.SelectBuilder, options *base.AccessOptions) (*ApiPager, error) {
	if options != nil && len(options.Roles) > 0 {
		value, _ := options.Roles.ToStringSlice()

		query = query.Where(squirrel.NewExprSlice(fmt.Sprintf("\"%s\" && ARRAY["+sq.Placeholders(len(options.Roles))+"]", "access"), value))
	}

	a.Manager.Remove(query)

	return a.Find(query, 0, 0, options)
}
