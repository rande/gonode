package core

import (
	"container/list"
	"encoding/json"
	sq "github.com/lann/squirrel"
	"github.com/twinj/uuid"
)

var (
	emptyUuid = GetReference(uuid.New([]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}))
	rootUuid  = GetReference(uuid.New([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
)

func InterfaceToJsonMessage(ntype string, data interface{}) json.RawMessage {
	v, err := json.Marshal(data)

	PanicOnError(err)

	return v
}

func GetEmptyReference() Reference {
	return emptyUuid
}

func GetRootReference() Reference {
	return rootUuid
}

type NodeManager interface {
	SelectBuilder() sq.SelectBuilder
	FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List
	FindOneBy(query sq.SelectBuilder) *Node
	Find(uuid Reference) *Node
	Remove(query sq.SelectBuilder) error
	RemoveOne(node *Node) (*Node, error)
	Save(node *Node) (*Node, error)
	Notify(channel string, payload string)
	NewNode(t string) *Node
}
