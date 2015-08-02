// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

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
	Prefix   string
}

func (m *PgNodeManager) SelectBuilder() sq.SelectBuilder {
	return sq.
		Select("id, uuid, type, name, revision, version, created_at, updated_at, set_uuid, parent_uuid, slug, created_by, updated_by, data, meta, deleted, enabled, source, status, weight").
		From(m.Prefix + "_nodes").
		PlaceholderFormat(sq.Dollar)
}

func (m *PgNodeManager) Notify(channel string, payload string) {
	m.Logger.Printf("[PgNode] NOTIFY %s, %s ", channel, payload)

	_, err := m.Db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", channel, strings.Replace(payload, "'", "''", -1)))

	PanicOnError(err)
}

func (m *PgNodeManager) NewNode(t string) *Node {
	return m.Handlers.NewNode(t)
}

func (m *PgNodeManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	query = query.Limit(limit).Offset(offset)

	if m.Logger != nil {
		rawSql, _, _ := query.ToSql()
		m.Logger.Print("[PgNode] FindBy: ", rawSql)
	}

	rows, err := query.
		RunWith(m.Db).
		Query()

	list := list.New()

	if err != nil {
		if m.Logger != nil {
			rawSql, _, _ := query.ToSql()
			m.Logger.Printf("[PgNode] Error while runing the request %s, %s ", rawSql, err)
		}

		return list
	}

	for rows.Next() {
		node := m.hydrate(rows)

		list.PushBack(node)
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
		&node.Version,
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

	PanicOnError(err)

	var tmpUuid uuid.UUID

	// transform UUID
	tmpUuid, _ = uuid.Parse(Uuid)
	node.Uuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.Parse(SetUuid)
	node.SetUuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.Parse(ParentUuid)
	node.ParentUuid = GetReference(tmpUuid)
	tmpUuid, _ = uuid.Parse(CreatedBy)
	node.CreatedBy = GetReference(tmpUuid)
	tmpUuid, _ = uuid.Parse(UpdatedBy)
	node.UpdatedBy = GetReference(tmpUuid)
	tmpUuid, _ = uuid.Parse(Source)
	node.Source = GetReference(tmpUuid)

	m.Handlers.Get(node).Load(data, meta, node)

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

			m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
				Type:     node.Type,
				Name:     node.Name,
				Action:   "SoftDelete",
				Subject:  node.Uuid.CleanString(),
				Revision: node.Revision,
				Date:     node.UpdatedAt,
			})

			m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)
		}
	}
}

func (m *PgNodeManager) RemoveOne(node *Node) (*Node, error) {
	node.UpdatedAt = time.Now()
	node.Deleted = true

	m.Logger.Printf("[PgNode] Soft Delete: Uuid:%+v - type: %s", node.Uuid, node.Type)

	m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
		Type:     node.Type,
		Action:   "SoftDelete",
		Subject:  node.Uuid.CleanString(),
		Revision: node.Revision,
		Date:     node.UpdatedAt,
		Name:     node.Name,
	})

	return m.Save(node)
}

func (m *PgNodeManager) insertNode(node *Node, table string) (*Node, error) {
	if node.Uuid == GetEmptyReference() {
		node.Uuid = GetReference(uuid.NewV4())
	}

	if node.Slug == "" {
		node.Slug = node.Uuid.String()
	}

	query := sq.Insert(table).
		Columns(
		"uuid", "type", "revision", "version", "name", "created_at", "updated_at", "set_uuid",
		"parent_uuid", "slug", "created_by", "updated_by", "data", "meta", "deleted",
		"enabled", "source", "status", "weight").
		Values(
		node.Uuid.CleanString(),
		node.Type,
		node.Revision,
		node.Version,
		node.Name,
		node.CreatedAt,
		node.UpdatedAt,
		node.SetUuid.CleanString(),
		node.ParentUuid.CleanString(),
		node.Slug,
		node.CreatedBy.CleanString(),
		node.UpdatedBy.CleanString(),
		string(InterfaceToJsonMessage(node.Type, node.Data)[:]),
		string(InterfaceToJsonMessage(node.Type, node.Meta)[:]),
		node.Deleted,
		node.Enabled,
		node.Source.CleanString(),
		node.Status,
		node.Weight).
		Suffix("RETURNING \"id\"").
		RunWith(m.Db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&node.id)

	return node, err
}

func (m *PgNodeManager) updateNode(node *Node, table string) (*Node, error) {
	var err error

	query := sq.Update(m.Prefix+"_nodes").RunWith(m.Db).PlaceholderFormat(sq.Dollar).
		Set("uuid", node.Uuid.CleanString()).
		Set("type", node.Type).
		Set("revision", node.Revision).
		Set("version", node.Version).
		Set("name", node.Name).
		Set("created_at", node.CreatedAt).
		Set("updated_at", node.UpdatedAt).
		Set("set_uuid", node.SetUuid.CleanString()).
		Set("parent_uuid", node.ParentUuid.CleanString()).
		Set("slug", node.Slug).
		Set("created_by", node.CreatedBy.CleanString()).
		Set("updated_by", node.UpdatedBy.CleanString()).
		Set("deleted", node.Deleted).
		Set("enabled", node.Enabled).
		Set("data", string(InterfaceToJsonMessage(node.Type, node.Data)[:])).
		Set("meta", string(InterfaceToJsonMessage(node.Type, node.Meta)[:])).
		Set("source", node.Source.CleanString()).
		Set("status", node.Status).
		Set("weight", node.Weight).
		Where("id = ?", node.id)

	_, err = query.Exec()

	PanicOnError(err)

	if m.Logger != nil {
		strQuery, _, _ := query.ToSql()
		m.Logger.Printf("[PgNode] Update: %s", strQuery)
	}

	return node, err
}

func (m *PgNodeManager) Save(node *Node) (*Node, error) {

	PanicIf(m.ReadOnly, "The manager is readonly, cannot alter the datastore")

	var err error

	handler := m.Handlers.Get(node)

	if node.id == 0 {
		handler.PreInsert(node, m)

		node, err = m.insertNode(node, m.Prefix+"_nodes")
		PanicOnError(err)
		node, err = m.insertNode(node, m.Prefix+"_nodes_audit")
		PanicOnError(err)

		if m.Logger != nil {
			m.Logger.Printf("[PgNode] Creating node uuid: %s, id: %d, type: %s", node.Uuid, node.id, node.Type)
		}

		handler.PostInsert(node, m)

		m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
			Type:    node.Type,
			Action:  "Create",
			Subject: node.Uuid.CleanString(),
			Date:    node.CreatedAt,
			Name:    node.Name,
		})

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
	saved, err = m.insertNode(saved, m.Prefix+"_nodes_audit")

	PanicOnError(err)

	// 3. Update the revision number
	node.Revision++
	node.CreatedAt = saved.CreatedAt
	node.UpdatedAt = saved.UpdatedAt

	node, err = m.updateNode(node, m.Prefix+"_nodes")

	handler.PostUpdate(node, m)

	PanicOnError(err)

	m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
		Type:     node.Type,
		Action:   "Update",
		Subject:  node.Uuid.CleanString(),
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
