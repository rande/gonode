package gonode

import (
	"io"
	"encoding/json"
	sq "github.com/lann/squirrel"
	"github.com/twinj/uuid"
)

type Pager struct {
	Elements []*Node
	Page     uint64
	Next     uint64
	Previous uint64
	Offset   uint64
	Limit    uint64
}

type Api struct {
	Version string
	Manager *PgNodeManager
}

func (a *Api) serialize(w io.Writer, data interface {}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)

	if err != nil {
		panic(err)
	}
}

func (a *Api) deserialize(r io.Reader, data interface {}) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(data)

	if err != nil {
		panic(err)
	}
}

func (a *Api) Find(w io.Writer, b sq.SelectBuilder, offset uint64, limit uint64) error {
	query := a.Manager.SelectBuilder()
	list := a.Manager.FindBy(query, offset, limit + 1)

	pager := &Pager{}
	pager.Elements = make([]*Node, limit)

	if (uint64(list.Len()) == limit + 1) {
		pager.Next = offset / limit + 1
	}

	if list.Len() > 0 {
		element := list.Front()
		for pos, _ := range pager.Elements {
			pager.Elements[pos] = element.Value.(*Node)

			element = element.Next()

			if element == nil {
				break
			}
		}
	}

	a.serialize(w, pager)

	return nil
}

func (a *Api) Save(r io.Reader, w io.Writer) error {
	node := NewNode()

	a.deserialize(r, node)

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		if node.Type != saved.Type {
			panic("Type mismatch")
		}

		if saved.Deleted {
			panic("Cannot save a deleted node, restore it first ...")
		}

		if node.Revision != saved.Revision {
			panic("Revision mismatch, please saved from the latest revision")
		}

		node.id = saved.id
	}

	a.Manager.Save(node)

	a.serialize(w, node)

	return nil
}

func (a *Api) FindOne(reference string, w io.Writer) error {
	v, err := uuid.ParseUUID(reference)

	if err != nil {
		panic(err)
	}

	a.serialize(w, a.Manager.Find(v))

	return nil
}

func (a *Api) RemoveOne(reference string, w io.Writer) error {
	v, err := uuid.ParseUUID(reference)

	if err != nil {
		panic(err)
	}

	node := a.Manager.Find(v)

	a.Manager.RemoveOne(node)

	a.serialize(w, node)

	return nil
}

func (a *Api) Remove(b sq.SelectBuilder, w io.Writer) error {
	a.Manager.Remove(b)

	a.Find(w, b, 0, 0)

	return nil
}
