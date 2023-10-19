// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"container/list"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/rande/gonode/core/helper"
)

var (
	emptyUuid = "1111111111111111"
	rootUuid  = "0000000000000000"
)

func InterfaceToJsonMessage(ntype string, data interface{}) json.RawMessage {
	v, err := json.Marshal(data)

	helper.PanicOnError(err)

	return v
}

func GetEmptyReference() string {
	return emptyUuid
}

func GetRootReference() string {
	return rootUuid
}

type NodeManager interface {
	SelectBuilder(option *SelectOptions) sq.SelectBuilder
	FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List
	FindOneBy(query sq.SelectBuilder) *Node
	Find(uuid string) *Node
	Remove(query sq.SelectBuilder) error
	RemoveOne(node *Node) (*Node, error)
	Save(node *Node, revision bool) (*Node, error)
	Notify(channel string, payload string)
	NewNode(t string) *Node
	Validate(node *Node) (bool, Errors)
	Move(uuid, parent string) (int64, error)
}
