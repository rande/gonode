// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"container/list"
	"encoding/json"

	sq "github.com/lann/squirrel"
	"github.com/rande/gonode/core/helper"
	"github.com/twinj/uuid"
)

var (
	emptyUuid = GetReference(uuid.New([]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}))
	rootUuid  = GetReference(uuid.New([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
)

func InterfaceToJsonMessage(ntype string, data interface{}) json.RawMessage {
	v, err := json.Marshal(data)

	helper.PanicOnError(err)

	return v
}

func GetEmptyReference() Reference {
	return emptyUuid
}

func GetRootReference() Reference {
	return rootUuid
}

type NodeManager interface {
	SelectBuilder(option *SelectOptions) sq.SelectBuilder
	FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List
	FindOneBy(query sq.SelectBuilder) *Node
	Find(uuid Reference) *Node
	Remove(query sq.SelectBuilder) error
	RemoveOne(node *Node) (*Node, error)
	Save(node *Node, revision bool) (*Node, error)
	Notify(channel string, payload string)
	NewNode(t string) *Node
	Validate(node *Node) (bool, Errors)
	Move(uuid, parent Reference) (int64, error)
}
