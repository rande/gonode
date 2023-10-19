// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/core/squirrel"
	"github.com/rande/gonode/modules/base"
	log "github.com/sirupsen/logrus"
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
		a.Logger.Printf("trying to save node.nid=%s, node.type=%s", node.Nid, node.Type)
	}

	saved := a.Manager.Find(node.Nid)

	if saved != nil {
		a.Logger.Printf("find nid: %s", node.Nid)

		helper.PanicUnless(node.Type == saved.Type, "Type mismatch")

		if options != nil {
			result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

			if !result {
				return nil, nil, base.ErrAccessForbidden
			}
		}

		if node.Deleted == true {
			return nil, nil, base.ErrAlreadyDeleted
		}

		if node.Revision != saved.Revision {
			return nil, nil, base.ErrRevision
		}

		node.Id = saved.Id

		// we cannot overwrite the Parents, Or the ParentNid, need to use the http API
		node.Parents = saved.Parents
		node.ParentNid = saved.ParentNid

	} else if a.Logger != nil {
		a.Logger.Printf("cannot find nid: %s, create a new one", node.Nid)
	}

	if a.Logger != nil {
		a.Logger.Printf("saving node.id=%d, node.nid=%s", node.Id, node.Nid)
	}

	if ok, errors := a.Manager.Validate(node); !ok {
		return nil, errors, base.ErrValidation
	}

	node, err := a.Manager.Save(node, true)

	return node, nil, err
}

func (a *Api) Move(nodeNid, parentNid string, options *base.AccessOptions) (*ApiOperation, error) {
	// handle node
	if node := a.Manager.Find(nodeNid); node == nil {
		return nil, base.ErrNotFound
	} else if result, _ := a.Authorizer.IsGranted(options.Token, nil, node); !result {
		return nil, base.ErrAccessForbidden
	}

	if parent := a.Manager.Find(parentNid); parent == nil {
		return nil, base.ErrNotFound
	} else if result, _ := a.Authorizer.IsGranted(options.Token, nil, parent); !result {
		return nil, base.ErrAccessForbidden
	}

	// move node
	if affectedNodes, err := a.Manager.Move(nodeNid, parentNid); err != nil {
		return nil, err
	} else {
		return &ApiOperation{
			Status:  OPERATION_OK,
			Message: fmt.Sprintf("Node altered: %d", affectedNodes),
		}, nil
	}
}

func (a *Api) FindOne(nid string, options *base.AccessOptions) (*base.Node, error) {
	query := a.Manager.SelectBuilder(base.NewSelectOptions()).Where(sq.Eq{"nid": nid})

	return a.FindOneBy(query, options)
}

func (a *Api) FindOneBy(query sq.SelectBuilder, options *base.AccessOptions) (*base.Node, error) {
	node := a.Manager.FindOneBy(query)

	if node == nil {
		return nil, base.ErrNotFound
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

		if !result {
			return nil, base.ErrAccessForbidden
		}
	}

	return node, nil
}

func (a *Api) RemoveOne(nid string, options *base.AccessOptions) (*base.Node, error) {
	node := a.Manager.Find(nid)

	if node == nil {
		return nil, base.ErrNotFound
	}

	if options != nil {
		result, _ := a.Authorizer.IsGranted(options.Token, nil, node)

		if !result {
			return nil, base.ErrAccessForbidden
		}
	}

	if node.Deleted {
		return nil, base.ErrAlreadyDeleted
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
