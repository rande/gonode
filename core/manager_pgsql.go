// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"errors"
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

type SelectOptions struct {
	TableSuffix  string
	SelectClause string
}

func NewSelectOptions() *SelectOptions {
	return &SelectOptions{
		TableSuffix:  "nodes",
		SelectClause: "id, uuid, type, name, revision, version, created_at, updated_at, set_uuid, parent_uuid, parents, slug, created_by, updated_by, data, meta, deleted, enabled, source, status, weight",
	}
}

func (m *PgNodeManager) SelectBuilder(options *SelectOptions) sq.SelectBuilder {
	return sq.
		Select(options.SelectClause).
		From(m.Prefix + "_" + options.TableSuffix).
		PlaceholderFormat(sq.Dollar)
}

func (m *PgNodeManager) Notify(channel string, payload string) {
	//	m.Logger.Printf("[PgNode] NOTIFY %s, %s ", channel, payload)

	_, err := m.Db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", channel, strings.Replace(payload, "'", "''", -1)))

	PanicOnError(err)
}

func (m *PgNodeManager) NewNode(t string) *Node {
	return m.Handlers.NewNode(t)
}

func (m *PgNodeManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	query = query.Limit(limit).Offset(offset)

	rawSql, _, _ := query.ToSql()

	rows, err := query.
		RunWith(m.Db).
		Query()

	list := list.New()

	if err != nil {
		if m.Logger != nil {
			m.Logger.Printf("[PgNode] Error while runing the request: `%s`, %s ", rawSql, err)
		}

		PanicOnError(err)
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
	return m.FindOneBy(m.SelectBuilder(NewSelectOptions()).Where(sq.Eq{"uuid": uuid.String()}))
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

	var Parents StringSlice

	err := rows.Scan(
		&node.Id,
		&Uuid,
		&node.Type,
		&node.Name,
		&node.Revision,
		&node.Version,
		&node.CreatedAt,
		&node.UpdatedAt,
		&SetUuid,
		&ParentUuid,
		&Parents,
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

	pUuids := make([]Reference, 0)

	for _, ref := range Parents {
		pUuid, _ := uuid.Parse(ref)
		pUuids = append(pUuids, GetReference(pUuid))
	}

	node.Parents = pUuids

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

			m.Save(node, false)

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

	return m.Save(node, true)
}

func (m *PgNodeManager) insertNode(node *Node, table string) (*Node, error) {
	if node.Uuid == GetEmptyReference() {
		node.Uuid = GetReference(uuid.NewV4())
	}

	if node.Slug == "" {
		node.Slug = node.Uuid.String()
	}

	Parents := make(StringSlice, 0)
	for _, p := range node.Parents {
		Parents = append(Parents, p.CleanString())
	}

	query := sq.Insert(table).
		Columns(
		"uuid", "type", "revision", "version", "name", "created_at", "updated_at", "set_uuid",
		"parent_uuid", "parents", "slug", "created_by", "updated_by", "data", "meta", "deleted",
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
		Parents,
		node.Slug,
		node.CreatedBy.CleanString(),
		node.UpdatedBy.CleanString(),
		string(InterfaceToJsonMessage(node.Type, node.Data)[:]),
		string(InterfaceToJsonMessage(node.Type, node.Meta)[:]),
		node.Deleted,
		node.Enabled,
		node.Source.CleanString(),
		node.Status,
		node.Weight,
	).
		Suffix("RETURNING \"id\"").
		RunWith(m.Db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&node.Id)

	return node, err
}

func (m *PgNodeManager) Move(uuid, parentUuid Reference) (int64, error) {

	tx, err := m.Db.Begin()

	if err != nil {
		return 0, err
	}

	r, err := tx.Exec(fmt.Sprintf(`UPDATE %s SET parent_uuid = $1 WHERE uuid = $2 AND EXISTS(SELECT uuid FROM %s WHERE uuid = $3 and $4 <> ALL(parents))`,
		m.Prefix+"_nodes", m.Prefix+"_nodes"),
		parentUuid.CleanString(),
		uuid.CleanString(),
		parentUuid.CleanString(),
		uuid.CleanString())

	if err != nil {
		tx.Rollback()

		return 0, err
	}

	affectedRows, err := r.RowsAffected()

	if err != nil {
		tx.Rollback()

		return 0, err
	}

	if affectedRows > 0 {
		tx.Exec(fmt.Sprintf(`WITH RECURSIVE  r AS (
					SELECT uuid, parent_uuid, parents
					FROM %s r
					WHERE uuid = $1::uuid
				UNION ALL
					SELECT c.uuid, c.parent_uuid, array_append(r.parents, c.parent_uuid) AS parents
					FROM %s c
					JOIN r ON c.parent_uuid = r.uuid
			)
			UPDATE %s n SET parents = r.parents FROM r WHERE r.uuid = n.uuid`,
			m.Prefix+"_nodes",
			m.Prefix+"_nodes",
			m.Prefix+"_nodes"),
			parentUuid.CleanString())
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()

		return 0, err
	}

	return affectedRows, nil
}

func (m *PgNodeManager) updateNode(node *Node, table string) (*Node, error) {

	PanicIf(node.Id == 0, "Cannot update node without id")

	query := sq.Update(m.Prefix+"_nodes").RunWith(m.Db).PlaceholderFormat(sq.Dollar).
		Set("uuid", node.Uuid.CleanString()).
		Set("type", node.Type).
		Set("revision", node.Revision).
		Set("version", node.Version).
		Set("name", node.Name).
		Set("created_at", node.CreatedAt).
		Set("updated_at", node.UpdatedAt).
		Set("set_uuid", node.SetUuid.CleanString()).
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
		Where("id = ?", node.Id)

	result, err := query.Exec()

	PanicOnError(err)

	affected, err := result.RowsAffected()

	PanicOnError(err)

	if affected == 0 {
		return node, errors.New("Zero affected rows for current node")
	}

	return node, err
}

func (m *PgNodeManager) Save(node *Node, revision bool) (*Node, error) {
	if m.Logger != nil {
		m.Logger.Printf("[PgNode] Saving uuid: %s, id: %d, type: %s, revision: %d", node.Uuid, node.Id, node.Type, node.Revision)
	}

	PanicIf(m.ReadOnly, "The manager is readonly, cannot alter the datastore")

	var err error

	handler := m.Handlers.Get(node)

	if node.Id == 0 {
		handler.PreInsert(node, m)

		node, err = m.insertNode(node, m.Prefix+"_nodes_audit")
		PanicOnError(err)

		node.Id = 0

		node, err = m.insertNode(node, m.Prefix+"_nodes")
		PanicOnError(err)

		if m.Logger != nil {
			m.Logger.Printf("[PgNode] Creating node uuid: %s, id: %d, type: %s, revision: %d", node.Uuid, node.Id, node.Type, node.Revision)
		}

		handler.PostInsert(node, m)

		m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
			Type:        node.Type,
			Action:      "Create",
			Subject:     node.Uuid.CleanString(),
			Date:        node.CreatedAt,
			Name:        node.Name,
			Revision:    node.Revision,
			NewRevision: revision,
		})

		return node, err
	}

	handler.PreUpdate(node, m)

	// 1. check if the one in the datastore is older
	saved := m.FindOneBy(m.SelectBuilder(NewSelectOptions()).Where(sq.Eq{"uuid": node.Uuid.String()}))

	if saved != nil && node.Revision != saved.Revision {
		m.Logger.Printf("[PgNode] Invalid revision for node: %s, saved rev: %d, current rev: %d", node.Uuid, saved.Revision, node.Revision)

		return node, NewRevisionError(fmt.Sprintf("Invalid revision for node: %s, saved rev: %d, current rev: %d", node.Uuid, saved.Revision, node.Revision))
	}

	if m.Logger != nil {
		m.Logger.Printf("[PgNode] Updating uuid: %s, id: %d, type: %s, revision: %d", node.Uuid, node.Id, node.Type, node.Revision)
	}

	if revision {
		// 3. Update the revision number
		node.Revision++
		node.CreatedAt = saved.CreatedAt
		node.UpdatedAt = saved.UpdatedAt

		m.Logger.Printf("[PgNode] Increment revision - uuid: %s, id: %d, type: %s, revision: %d", node.Uuid, node.Id, node.Type, node.Revision)
	}

	node, err = m.updateNode(node, m.Prefix+"_nodes")
	PanicOnError(err)

	handler.PostUpdate(node, m)

	if revision {
		id := node.Id
		_, err = m.insertNode(node, m.Prefix+"_nodes_audit")

		node.Id = id
		PanicOnError(err)
	}

	m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
		Type:        node.Type,
		Action:      "Update",
		Subject:     node.Uuid.CleanString(),
		Revision:    node.Revision,
		Date:        node.UpdatedAt,
		Name:        node.Name,
		NewRevision: revision,
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
		errors.AddError("name", "Username cannot be empty")
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
