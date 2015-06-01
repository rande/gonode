package core

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/lann/squirrel"
	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
	"log"
	"strings"
	"time"
)

type PgNodeManager struct {
	Logger   *log.Logger
	Handlers Handlers
	Db       *sql.DB
	ReadOnly bool
}

func (m *PgNodeManager) SelectBuilder() sq.SelectBuilder {
	return sq.
		Select("id, uuid, type, name, revision, created_at, updated_at, set_uuid, parent_uuid, slug, created_by, updated_by, data, meta, deleted, enabled, source, status, weight").
		From("nodes").
		PlaceholderFormat(sq.Dollar)
}

func (m *PgNodeManager) Notify(channel string, payload string) {
	m.Logger.Printf("[PgNode] NOTIFY %s, %s ", channel, payload)

	_, err := m.Db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", channel, strings.Replace(payload, "'", "''", -1)))

	if err != nil {
		panic(err)
	}
}

func (m *PgNodeManager) NewNode(t string) *Node {
	return m.Handlers.NewNode(t)
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

func (m *PgNodeManager) Find(uuid Reference) *Node {
	return m.FindOneBy(m.SelectBuilder().Where(sq.Eq{"uuid": uuid.String()}))
}

func (m *PgNodeManager) hydrate(rows *sql.Rows) *Node {
	node := &Node{}

	data := json.RawMessage{}
	meta := json.RawMessage{}

	Uuid := ""
	SetUuid := ""
	ParentUuid := ""
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
		&node.Enabled,
		&Source,
		&node.Status,
		&node.Weight,
	)

	var tmpUuid uuid.UUID

	// transform UUID
	tmpUuid, _ = uuid.ParseUUID(Uuid)
	node.Uuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.ParseUUID(SetUuid)
	node.SetUuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.ParseUUID(ParentUuid)
	node.ParentUuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.ParseUUID(CreatedBy)
	node.CreatedBy = GetReference(tmpUuid)
	tmpUuid, _ = uuid.ParseUUID(UpdatedBy)
	node.UpdatedBy = GetReference(tmpUuid)
	tmpUuid, _ = uuid.ParseUUID(Source)
	node.Source = GetReference(tmpUuid)

	node.Data, node.Meta = m.Handlers.Get(node).GetStruct()

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

func (m *PgNodeManager) Remove(query sq.SelectBuilder) error {
	query = query.Where("deleted != ?", true)

	now := time.Now()

	for {
		nodes := m.FindBy(query, 0, 1024)

		if nodes.Len() == 0 {
			return nil
		}

		for e := nodes.Front(); e != nil; e = e.Next() {
			node := e.Value.(*Node)
			node.Deleted = true
			node.UpdatedAt = now

			m.Save(node)

			m.sendNotification("manager_action", &ModelEvent{
				Type:     node.Type,
				Name:     node.Name,
				Action:   "SoftDelete",
				Subject:  uuid.Formatter(node.Uuid, uuid.CleanHyphen),
				Revision: node.Revision,
				Date:     node.UpdatedAt,
			})

			m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)
		}
	}

	return nil
}

func (m *PgNodeManager) RemoveOne(node *Node) (*Node, error) {
	node.UpdatedAt = time.Now()
	node.Deleted = true

	m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)

	m.sendNotification("manager_action", &ModelEvent{
		Type:     node.Type,
		Action:   "SoftDelete",
		Subject:  uuid.Formatter(node.Uuid, uuid.CleanHyphen),
		Revision: node.Revision,
		Date:     node.UpdatedAt,
		Name:     node.Name,
	})

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
	m.Logger.Printf("[PgNode]  > Status:     %s", node.Status)
	m.Logger.Printf("[PgNode]  > Weight:     %s", node.Weight)
	m.Logger.Printf("[PgNode]  > Deleted:    %s", node.Deleted)
	m.Logger.Printf("[PgNode]  > Enabled:    %s", node.Enabled)
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

	if node.Uuid == GetEmptyReference() {
		node.Uuid = GetReference(uuid.NewV4())
	}

	if node.Slug == "" {
		node.Slug = node.Uuid.String()
	}

	query := sq.Insert(table).
		Columns(
		"uuid", "type", "revision", "name", "created_at", "updated_at", "set_uuid",
		"parent_uuid", "slug", "created_by", "updated_by", "data", "meta", "deleted",
		"enabled", "source", "status", "weight").
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
		node.Enabled,
		uuid.Formatter(node.Source, uuid.CleanHyphen),
		node.Status,
		node.Weight).
		Suffix("RETURNING \"id\"").
		RunWith(m.Db).
		PlaceholderFormat(sq.Dollar)

	query.QueryRow().Scan(&node.id)

	return node, err
}

func (m *PgNodeManager) updateNode(node *Node, table string) (*Node, error) {
	var err error

	query := sq.Update("nodes").RunWith(m.Db).PlaceholderFormat(sq.Dollar).
		Set("uuid", uuid.Formatter(node.Uuid, uuid.CleanHyphen)).
		Set("type", node.Type).
		Set("revision", node.Revision).
		Set("name", node.Name).
		Set("created_at", node.CreatedAt).
		Set("updated_at", node.UpdatedAt).
		Set("set_uuid", uuid.Formatter(node.SetUuid, uuid.CleanHyphen)).
		Set("parent_uuid", uuid.Formatter(node.ParentUuid, uuid.CleanHyphen)).
		Set("slug", node.Slug).
		Set("created_by", uuid.Formatter(node.CreatedBy, uuid.CleanHyphen)).
		Set("updated_by", uuid.Formatter(node.UpdatedBy, uuid.CleanHyphen)).
		Set("deleted", node.Deleted).
		Set("enabled", node.Enabled).
		Set("data", string(InterfaceToJsonMessage(node.Type, node.Data)[:])).
		Set("meta", string(InterfaceToJsonMessage(node.Type, node.Meta)[:])).
		Set("source", uuid.Formatter(node.Source, uuid.CleanHyphen)).
		Set("status", node.Status).
		Set("weight", node.Weight).
		Where("id = ?", node.id)

	_, err = query.Exec()

	if err != nil {
		log.Fatal(err)
	}

	if m.Logger != nil {
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
	handler := m.Handlers.Get(node)

	if node.id == 0 {
		handler.PreInsert(node, m)

		node, err = m.insertNode(node, "nodes")
		node, err = m.insertNode(node, "nodes_audit")

		if m.Logger != nil {
			m.Logger.Printf("[PgNode] Creating node uuid: %s, id: %d, type: %s", node.Uuid, node.id, node.Type)
		}

		m.sendNotification("manager_action", &ModelEvent{
			Type:    node.Type,
			Action:  "Create",
			Subject: uuid.Formatter(node.Uuid, uuid.CleanHyphen),
			Date:    node.CreatedAt,
			Name:    node.Name,
		})

		handler.PostInsert(node, m)

		return node, err
	}

	handler.PreUpdate(node, m)

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

	handler.PostUpdate(node, m)

	if err != nil {
		panic(err)
	}

	m.sendNotification("manager_action", &ModelEvent{
		Type:     node.Type,
		Action:   "Update",
		Subject:  uuid.Formatter(node.Uuid, uuid.CleanHyphen),
		Revision: node.Revision,
		Date:     node.UpdatedAt,
		Name:     node.Name,
	})

	return node, err
}

func (m *PgNodeManager) sendNotification(channel string, element interface{}) {
	data, _ := json.Marshal(element)

	m.Notify(channel, string(data[:]))
}

func (m *PgNodeManager) Validate(node *Node) (bool, Errors) {
	errors := NewErrors()

	if node.Name == "" {
		errors.AddError("name", "Login cannot be empty")
	}

	if node.Slug == "" {
		errors.AddError("slug", "Name cannot be empty")
	}

	if node.Type == "" {
		errors.AddError("type", "Type cannot be empty")
	}

	if node.Status < 0 || node.Status > 3 {
		errors.AddError("status", "Invalid status")
	}

	m.Handlers.Get(node).Validate(node, m, errors)

	return !errors.HasErrors(), errors
}
