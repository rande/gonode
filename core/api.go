package gonode

import (
	"io"
	"encoding/json"
	sq "github.com/lann/squirrel"
	"bytes"
)


type SearchForm struct {
	Page       uint64             `schema:"page"`
	PerPage    uint64             `schema:"per_page"`
	Uuid       string             `schema:"uuid"`
	Type       []string           `schema:"type"`
	Name       string             `schema:"name"`
	Slug       string             `schema:"slug"`
	Data       map[string][]string  `schema:"data"`
	Meta       map[string][]string  `schema:"meta"`
	Status     []string             `schema:"status"`
	Weight     []string             `schema:"weight"`
	Revision   string             `schema:"revision"`
//	CreatedAt  time.Time          `schema:"created_at"`
//	UpdatedAt  time.Time          `schema:"updated_at"`
	Enabled    string             `schema:"enabled"`
	Deleted    string             `schema:"deleted"`
	Current    string             `schema:"current"`
//	Parents    []Reference        `schema:"parents"`
	UpdatedBy  []string           `schema:"updated_by"`
	CreatedBy  []string           `schema:"created_by"`
	ParentUuid []string           `schema:"parent_uuid"`
	SetUuid    []string           `schema:"set_uuid"`
	Source     []string           `schema:"source"`
}

func GetSearchForm() *SearchForm {
	return &SearchForm{
		Data: make(map[string][]string),
		Meta: make(map[string][]string),
	}
}

type ApiPager struct {
	Elements []*Node `json:"elements"`
	Page     uint64  `json:"page"`
	PerPage  uint64  `json:"per_page"`
	Next     uint64  `json:"next"`
	Previous uint64  `json:"previous"`
}

type Api struct {
	Version string
	Manager *PgNodeManager
	BaseUrl string
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

func (a *Api) SelectBuilder() sq.SelectBuilder {
	return a.Manager.SelectBuilder()
}

func (a *Api) Find(w io.Writer, query sq.SelectBuilder, page uint64, perPage uint64) error {

	list := a.Manager.FindBy(query, page * perPage, perPage + 1)

	pager := &ApiPager{
		Page: page,
		PerPage: perPage,
	}

	if (uint64(list.Len()) == perPage + 1) {
		pager.Next = page + 1
		pager.Elements = make([]*Node, perPage)
	} else {
		pager.Elements = make([]*Node, list.Len())
	}

	if (page > 1) {
		pager.Previous = page - 1
	}

	a.Manager.Logger.Printf("Result len: %s", list.Len())

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

	// we need to deserialize twice to load the correct Meta/Data structure
	var data bytes.Buffer
	read, err := data.ReadFrom(r)

	reader := bytes.NewReader(data.Bytes())

	if err != nil {
		panic(err)
	}

	if (read == 0) {
		panic("no data read from the request")
	}

	a.deserialize(reader, node)

	reader.Seek(0, 0)

	node.Data, node.Meta = a.Manager.GetHandler(node).GetStruct()
	a.deserialize(reader, node)

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

func (a *Api) FindOne(uuid string, w io.Writer) error {

	a.serialize(w, a.Manager.Find(GetReferenceFromString(uuid)))

	return nil
}

func (a *Api) RemoveOne(uuid string, w io.Writer) error {

	node := a.Manager.Find(GetReferenceFromString(uuid))

	a.Manager.DumpNode(node)

	node, _ = a.Manager.RemoveOne(node)

	a.serialize(w, node)

	return nil
}

func (a *Api) Remove(b sq.SelectBuilder, w io.Writer) error {
	a.Manager.Remove(b)

	a.Find(w, b, 0, 0)

	return nil
}
