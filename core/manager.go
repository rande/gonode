package gonode

import (
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"encoding/json"
	"time"
	"github.com/twinj/uuid"
	"container/list"
	sq "github.com/lann/squirrel"
	"fmt"
)

var (
	emptyUuid       = uuid.New([]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11})
	rootUuid        = uuid.New([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	StatusDraft     = 1
	StatusCompleted = 2
	StatusValidated = 3
)

func InterfaceToJsonMessage(ntype string, data interface {}) json.RawMessage {
	v, _ := json.Marshal(data)

	return v
}

func GetEmptyUuid() uuid.UUID {
	return emptyUuid
}

func GetRootUuid() uuid.UUID {
	return rootUuid
}

type NodeManager interface {
	Find(query interface{}, offset int, limit int) []*Node
	FindOne(query interface{}) *Node
	Save(node *Node) (*Node, error)
	Remove(query interface{}) error
	RemoveOne(node *Node) (*Node, error)
}

type PgNodeManager struct {
	Logger     *log.Logger
	Handlers   map[string] interface {}
	Db         *sql.DB
	ReadOnly   bool
}

func (m *PgNodeManager) SelectBuilder() sq.SelectBuilder {
	return sq.
		Select("id, uuid, type, name, revision, created_at, updated_at, set_uuid, parent_uuid, slug, created_by, updated_by, data, meta, deleted, source").
		From("nodes").
		PlaceholderFormat(sq.Dollar)
}

func (m *PgNodeManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	query = query.Limit(limit).Offset(offset)

	if m.Logger != nil {
		sql, _, _ := query.ToSql()
		m.Logger.Print("[PgNode] FindBy: ", sql)
	}

	rows, err := query.
		RunWith(m.Db).
		Query()

	if err != nil {
		log.Fatal(err)
	}

	list := list.New()

	for rows.Next() {
		list.PushBack(m.hydrate(rows))
	}

	return list
}

func (m *PgNodeManager) FindOneBy(query sq.SelectBuilder) *Node {
	list := m.FindBy(query, 0, 1)

	if list.Len() == 1 {
		return list.Front().Value.(*Node)
	}

	return nil
}

func (m *PgNodeManager) Find(uuid uuid.UUID) *Node {
	return m.FindOneBy(m.SelectBuilder().Where(sq.Eq{"uuid": uuid.String(), "deleted": false}))
}

func (m *PgNodeManager) hydrate(rows *sql.Rows) *Node {
	node := &Node{}

	data := json.RawMessage{}
	meta := json.RawMessage{}

	Uuid := ""
	SetUuid := ""
	ParentUuid :=  ""
	CreatedBy := ""
	UpdatedBy := ""
	Source := ""

	err := rows.Scan(
		&node.id,
		&Uuid,
		&node.Type,
		&node.Name,
		&node.Revision,
		&node.CreatedAt,
		&node.UpdatedAt,
		&SetUuid,
		&ParentUuid,
		&node.Slug,
		&CreatedBy,
		&UpdatedBy,
		&data,
		&meta,
		&node.Deleted,
		&Source,
	)

	// transform UUID
	node.Uuid, _       = uuid.ParseUUID(Uuid)
	node.SetUuid, _    = uuid.ParseUUID(SetUuid)
	node.CreatedBy, _  = uuid.ParseUUID(ParentUuid)
	node.UpdatedBy, _  = uuid.ParseUUID(CreatedBy)
	node.ParentUuid, _ = uuid.ParseUUID(UpdatedBy)
	node.Source, _     = uuid.ParseUUID(Source)

	node.Data, node.Meta = m.Handlers[node.Type].(Handler).GetStruct()

	err = json.Unmarshal(data, node.Data)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(meta, node.Meta)
	if err != nil {
		log.Fatal(err)
	}

	return node
}

func (m *PgNodeManager) Remove(query sq.SelectBuilder) (error) {
	query = query.Where("deleted != ?", true)

	for {
		nodes := m.FindBy(query, 0, 1024)

		if nodes.Len() == 0 {
			return nil
		}

		for e := nodes.Front(); e != nil; e = e.Next() {
			node := e.Value.(*Node)
			node.Deleted = true
			m.Save(node)

			m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)
		}
	}

	return nil
}

func (m *PgNodeManager) RemoveOne(node *Node) (*Node, error) {
	node.Deleted = true

	m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)

	return m.Save(node)
}

func (m *PgNodeManager) DumpNode(node *Node) {
	if node == nil {
		panic("Cannot dump, node is nil")
	}

	m.Logger.Printf("[PgNode] ---- Node: %+v", node.id)
	m.Logger.Printf("[PgNode]  > Uuid:       %s", node.Uuid)
	m.Logger.Printf("[PgNode]  > Type:       %s", node.Type)
	m.Logger.Printf("[PgNode]  > Name:       %s", node.Name)
	m.Logger.Printf("[PgNode]  > Revision:   %d", node.Revision)
	m.Logger.Printf("[PgNode]  > CreatedAt:  %+v", node.CreatedAt)
	m.Logger.Printf("[PgNode]  > UpdatedAt:  %+v", node.UpdatedAt)
	m.Logger.Printf("[PgNode]  > Slug:       %s", node.Slug)
	m.Logger.Printf("[PgNode]  > Data:       %T => %+v", node.Data, node.Data)
	m.Logger.Printf("[PgNode]  > Meta:       %T => %+v", node.Meta, node.Meta)
	m.Logger.Printf("[PgNode]  > CreatedBy:  %s", node.CreatedBy)
	m.Logger.Printf("[PgNode]  > UpdatedBy:  %s", node.UpdatedBy)
	m.Logger.Printf("[PgNode]  > ParentUuid: %s", node.ParentUuid)
	m.Logger.Printf("[PgNode]  > SetUuid:    %s", node.SetUuid)
	m.Logger.Printf("[PgNode]  > Source:     %s", node.Source)
	m.Logger.Printf("[PgNode] ---- End Node")
}

func (m *PgNodeManager) insertNode(node *Node, table string) (*Node, error) {
	var err error

	if node.Uuid == GetEmptyUuid() {
		node.Uuid = uuid.NewV4()
	}

	query := sq.Insert(table).
		Columns("uuid", "type", "revision", "name", "created_at", "updated_at", "set_uuid", "parent_uuid", "slug", "created_by", "updated_by", "data", "meta", "deleted", "source").
    	Values(
			uuid.Formatter(node.Uuid, uuid.CleanHyphen),
			node.Type,
			node.Revision,
			node.Name,
			node.CreatedAt,
			node.UpdatedAt,
			uuid.Formatter(node.SetUuid, uuid.CleanHyphen),
			uuid.Formatter(node.ParentUuid, uuid.CleanHyphen),
			node.Slug,
			uuid.Formatter(node.CreatedBy, uuid.CleanHyphen),
			uuid.Formatter(node.UpdatedBy, uuid.CleanHyphen),
			string(InterfaceToJsonMessage(node.Type, node.Data)[:]),
			string(InterfaceToJsonMessage(node.Type, node.Meta)[:]),
			node.Deleted,
			uuid.Formatter(node.Source, uuid.CleanHyphen)).
		Suffix("RETURNING \"id\"").
		RunWith(m.Db).
		PlaceholderFormat(sq.Dollar)

	query.QueryRow().Scan(&node.id)

	return node, err
}

func (m *PgNodeManager) updateNode(node *Node, table string) (*Node, error) {
	var err error

	query := sq.Update("nodes").RunWith(m.Db).PlaceholderFormat(sq.Dollar).
		Set("uuid",        uuid.Formatter(node.Uuid, uuid.CleanHyphen)).
		Set("type",        node.Type).
		Set("revision",    node.Revision).
		Set("name",        node.Name).
		Set("created_at",  node.CreatedAt).
		Set("updated_at",  node.UpdatedAt).
		Set("set_uuid",    uuid.Formatter(node.SetUuid, uuid.CleanHyphen)).
		Set("parent_uuid", uuid.Formatter(node.ParentUuid, uuid.CleanHyphen)).
		Set("slug",        node.Slug).
		Set("created_by",  uuid.Formatter(node.CreatedBy, uuid.CleanHyphen)).
		Set("updated_by",  uuid.Formatter(node.UpdatedBy, uuid.CleanHyphen)).
		Set("deleted",     node.Deleted).
		Set("data",        string(InterfaceToJsonMessage(node.Type, node.Data)[:])).
		Set("meta",        string(InterfaceToJsonMessage(node.Type, node.Meta)[:])).
		Set("source",      uuid.Formatter(node.Source, uuid.CleanHyphen)).

		Where("id = ?", node.id)

	_, err = query.Exec()

	if err != nil {
		log.Fatal(err)
	}

	if (m.Logger != nil) {
		strQuery, _, _ := query.ToSql()
		m.Logger.Printf("[PgNode] Update: %s", strQuery)
	}

	return node, err
}

func (m *PgNodeManager) Save(node *Node) (*Node, error) {
	if m.ReadOnly {
		panic("The manager is readonly, cannot alter the datastore")
	}

	var err error

	if node.id == 0 {
		node, err = m.insertNode(node, "nodes")
		node, err = m.insertNode(node, "nodes_audit")

		if m.Logger != nil {
			m.Logger.Printf("[PgNode] Creating node uuid: %s, id: %d, type: %s", node.Uuid, node.id, node.Type)
		}

		return node, err
	}

	if m.Logger != nil {
		m.Logger.Printf("[PgNode] Updating node uuid: %s with id: %d, type: %s", node.Uuid, node.id, node.Type)
	}

	// 1. check if the one in the datastore is older
	saved := m.FindOneBy(m.SelectBuilder().Where(sq.Eq{"uuid": node.Uuid.String()}))

	if saved != nil && node.Revision != saved.Revision {
		return node, NewRevisionError(fmt.Sprintf("Invalid revision for node:%s, current rev: %d", node.Uuid, node.Revision))
	}

	// 2. Flag the current node as deprecated
	saved.UpdatedAt = time.Now()
	saved, err = m.insertNode(saved, "nodes_audit")

	if err != nil {
		panic(err)
	}

	// 3. Update the revision number
	node.Revision++
	node.CreatedAt = saved.CreatedAt
	node.UpdatedAt = saved.UpdatedAt

	node, err = m.updateNode(node, "nodes")

	if err != nil {
		panic(err)
	}

	return node, err
}
