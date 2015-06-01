package core

import (
	sq "github.com/lann/squirrel"
	"io"
)

type SearchForm struct {
	Page     uint64              `schema:"page"`
	PerPage  uint64              `schema:"per_page"`
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
	Elements []*Node `json:"elements"`
	Page     uint64  `json:"page"`
	PerPage  uint64  `json:"per_page"`
	Next     uint64  `json:"next"`
	Previous uint64  `json:"previous"`
}

type Api struct {
	Version    string
	Manager    *PgNodeManager
	BaseUrl    string
	Serializer *Serializer
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

	if uint64(list.Len()) == perPage+1 {
		pager.Next = page + 1
		pager.Elements = make([]*Node, perPage)
	} else {
		pager.Elements = make([]*Node, list.Len())
	}

	if page > 1 {
		pager.Previous = page - 1
	}

	a.Manager.Logger.Printf("Result len: %s", list.Len())

	if list.Len() > 0 {
		element := list.Front()
		for pos := range pager.Elements {
			pager.Elements[pos] = element.Value.(*Node)

			element = element.Next()

			if element == nil {
				break
			}
		}
	}

	Serialize(w, pager)

	return nil
}

func (a *Api) Save(r io.Reader, w io.Writer) error {
	node := NewNode()

	a.Serializer.Deserialize(r, node)

	a.Manager.Logger.Printf("trying to save node.uuid=%s", node.Uuid)

	saved := a.Manager.Find(node.Uuid)

	if saved != nil {
		a.Manager.Logger.Printf("find uuid: %s", node.Uuid)

		if node.Type != saved.Type {
			panic("Type mismatch")
		}

		if saved.Deleted {
			panic("Cannot save a deleted node, restore it first ...")
		}

		if node.Revision != saved.Revision {
			return RevisionError
		}

		node.id = saved.id
	} else {
		a.Manager.Logger.Printf("cannot find uuid: %s", node.Uuid)
	}

	a.Manager.Logger.Printf("saving node.id=%s, node.uuid=%s", node.id, node.Uuid)

	if ok, errors := a.Manager.Validate(node); !ok {
		Serialize(w, errors)

		return ValidationError
	}

	a.Manager.Save(node)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) FindOne(uuid string, w io.Writer) error {
	a.Serializer.Serialize(w, a.Manager.Find(GetReferenceFromString(uuid)))

	return nil
}

func (a *Api) RemoveOne(uuid string, w io.Writer) error {
	node := a.Manager.Find(GetReferenceFromString(uuid))

	a.Manager.DumpNode(node)

	node, _ = a.Manager.RemoveOne(node)

	a.Serializer.Serialize(w, node)

	return nil
}

func (a *Api) Remove(b sq.SelectBuilder, w io.Writer) error {
	a.Manager.Remove(b)

	a.Find(w, b, 0, 0)

	return nil
}
