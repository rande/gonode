// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	sq "github.com/lann/squirrel"
	_ "github.com/lib/pq"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/core/squirrel"
	"github.com/twinj/uuid"
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
		SelectClause: "id, uuid, type, name, revision, version, created_at, updated_at, set_uuid, parent_uuid, parents, slug, path, created_by, updated_by, data, meta, modules, access, deleted, enabled, source, status, weight",
	}
}

func (m *PgNodeManager) SelectBuilder(options *SelectOptions) sq.SelectBuilder {
	if options == nil {
		options = NewSelectOptions()
	}

	return sq.
		Select(options.SelectClause).
		From(m.Prefix + "_" + options.TableSuffix).
		PlaceholderFormat(sq.Dollar)
}

func (m *PgNodeManager) Notify(channel string, payload string) {
	_, err := m.Db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", channel, strings.Replace(payload, "'", "''", -1)))

	helper.PanicOnError(err)
}

func (m *PgNodeManager) NewNode(t string) *Node {
	return m.Handlers.NewNode(t)
}

func (m *PgNodeManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	query = query.Limit(limit).Offset(offset)

	rows, err := query.
		RunWith(m.Db).
		Query()

	list := list.New()

	rawSql, _, _ := query.ToSql()

	if m.Logger != nil {
		m.Logger.WithFields(log.Fields{
			"module": "node.manager",
			"query":  rawSql,
		}).Debug("Executing query")
	}

	if err != nil {
		if m.Logger != nil {
			m.Logger.WithFields(log.Fields{
				"module": "node.manager",
				"err":    err,
				"query":  rawSql,
			}).Warn("error while running the query")
		}

		helper.PanicOnError(err)
	}

	for rows.Next() {
		//if m.Logger != nil {
		//	m.Logger.WithFields(log.Fields{
		//		"module": "node.manager",
		//		"row":  rows,
		//	}).Debug("hydrating row")
		//}

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
	modules := json.RawMessage{}

	Uuid := ""
	SetUuid := ""
	ParentUuid := ""
	CreatedBy := ""
	UpdatedBy := ""
	Source := ""

	var Access, Parents squirrel.StringSlice

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
		&node.Path,
		&CreatedBy,
		&UpdatedBy,
		&data,
		&meta,
		&modules,
		&Access,
		&node.Deleted,
		&node.Enabled,
		&Source,
		&node.Status,
		&node.Weight,
	)

	helper.PanicOnError(err)

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

	for _, access := range Access {
		node.Access = append(node.Access, access)
	}

	handler := m.Handlers.Get(node)
	if h, ok := handler.(LoadNodeHandler); ok {
		h.Load(data, meta, node)
	} else {
		HandlerLoad(handler, data, meta, node)
	}

	err = json.Unmarshal(modules, &node.Modules)
	helper.PanicOnError(err)

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

			if m.Logger != nil {
				m.Logger.WithFields(log.Fields{
					"type":   node.Type,
					"uuid":   node.Uuid,
					"module": "node.manager",
				}).Warn("soft delete many")
			}
		}
	}
}

func (m *PgNodeManager) RemoveOne(node *Node) (*Node, error) {
	node.UpdatedAt = time.Now()
	node.Deleted = true

	if m.Logger != nil {
		m.Logger.WithFields(log.Fields{
			"type":   node.Type,
			"uuid":   node.Uuid,
			"module": "node.manager",
		}).Warn("soft delete one")
	}

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
	if node.Uuid.String() == GetEmptyReference().String() {
		node.Uuid = GetReference(uuid.NewV4())
	}

	if node.Slug == "" {
		node.Slug = node.Uuid.String()
	}

	Parents := make(squirrel.StringSlice, 0)
	for _, p := range node.Parents {
		Parents = append(Parents, p.CleanString())
	}

	node.Access = security.EnsureRoles(node.Access, "node:api:master")

	Access := make(squirrel.StringSlice, 0)
	for _, a := range node.Access {
		Access = append(Access, a)
	}

	query := sq.Insert(table).
		Columns(
			"uuid", "type", "revision", "version", "name", "created_at", "updated_at", "set_uuid",
			"parent_uuid", "parents", "slug", "path", "created_by", "updated_by", "data", "meta", "modules",
			"access", "deleted", "enabled", "source", "status", "weight").
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
			node.Path,
			node.CreatedBy.CleanString(),
			node.UpdatedBy.CleanString(),
			string(InterfaceToJsonMessage(node.Type, node.Data)[:]),
			string(InterfaceToJsonMessage(node.Type, node.Meta)[:]),
			string(InterfaceToJsonMessage(node.Type, node.Modules)[:]),
			Access,
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
		// @todo: optimize to only use the 2 tree's branches: source and target
		//        to avoid rebuilding the full tree
		tx.Exec(fmt.Sprintf(`WITH RECURSIVE  r AS (
					SELECT uuid, parent_uuid, parents,
						CASE	WHEN type = 'core.root' THEN ''
							WHEN array_length(parents, 1)>0 THEN path
							ELSE '/' || slug::varchar(2000)
						END as path
					FROM %s r
					WHERE uuid = $1::uuid
				UNION ALL
					SELECT c.uuid, c.parent_uuid, array_append(r.parents, c.parent_uuid) AS parents, r.path || '/' || c.slug as path
					FROM %s c
					JOIN r ON c.parent_uuid = r.uuid
			)
			UPDATE %s n SET parents = r.parents, path = r.path FROM r WHERE r.uuid = n.uuid`,
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
	helper.PanicIf(node.Id == 0, "Cannot update node without id")

	node.Access = security.EnsureRoles(node.Access, "node:api:master")

	Access := make(squirrel.StringSlice, 0)
	for _, a := range node.Access {
		Access = append(Access, a)
	}

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
		Set("path", node.Path).
		Set("created_by", node.CreatedBy.CleanString()).
		Set("updated_by", node.UpdatedBy.CleanString()).
		Set("deleted", node.Deleted).
		Set("enabled", node.Enabled).
		Set("data", string(InterfaceToJsonMessage(node.Type, node.Data)[:])).
		Set("meta", string(InterfaceToJsonMessage(node.Type, node.Meta)[:])).
		Set("modules", string(InterfaceToJsonMessage(node.Type, node.Modules)[:])).
		Set("access", Access).
		Set("source", node.Source.CleanString()).
		Set("status", node.Status).
		Set("weight", node.Weight).
		Where("id = ?", node.Id)

	result, err := query.Exec()

	helper.PanicOnError(err)

	affected, err := result.RowsAffected()

	helper.PanicOnError(err)

	if affected == 0 {
		return node, errors.New("Zero affected rows for current node")
	}

	return node, err
}

func (m *PgNodeManager) Save(node *Node, revision bool) (*Node, error) {

	var contextLogger *log.Entry

	if m.Logger != nil {
		contextLogger = m.Logger.WithFields(log.Fields{
			"uuid":     node.Uuid,
			"id":       node.Id,
			"type":     node.Type,
			"revision": node.Revision,
			"module":   "node.manager",
		})

		contextLogger.Debug("saving node")
	}

	helper.PanicIf(m.ReadOnly, "The manager is readonly, cannot alter the datastore")

	var err error

	handler := m.Handlers.Get(node)

	if node.Id == 0 {

		if h, ok := handler.(DatabaseNodeHandler); ok {
			h.PreInsert(node, m)
		}

		node, err = m.insertNode(node, m.Prefix+"_nodes_audit")
		helper.PanicOnError(err)

		node.Id = 0

		node, err = m.insertNode(node, m.Prefix+"_nodes")
		helper.PanicOnError(err)

		if contextLogger != nil {
			contextLogger.Debug("creating node")
		}

		if h, ok := handler.(DatabaseNodeHandler); ok {
			h.PostInsert(node, m)
		}

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

	if h, ok := handler.(DatabaseNodeHandler); ok {
		h.PreUpdate(node, m)
	}

	// 1. check if the one in the datastore is older
	saved := m.FindOneBy(m.SelectBuilder(NewSelectOptions()).Where(sq.Eq{"uuid": node.Uuid.String()}))

	if saved != nil && node.Revision != saved.Revision {
		if contextLogger != nil {
			contextLogger.Info("invalid revision for node")
		}

		return node, NewRevisionError(fmt.Sprintf("Invalid revision for node: %s, saved rev: %d, current rev: %d", node.Uuid, saved.Revision, node.Revision))
	}

	if contextLogger != nil {
		contextLogger.Debug("updating node")
	}

	if revision {
		// 3. Update the revision number
		node.Revision++
		node.CreatedAt = saved.CreatedAt

		if contextLogger != nil {
			contextLogger.Debug("increment revision")
		}
	}

	updatedAt := time.Now()

	node.UpdatedAt = updatedAt

	node, err = m.updateNode(node, m.Prefix+"_nodes")
	helper.PanicOnError(err)

	if h, ok := handler.(DatabaseNodeHandler); ok {
		h.PostUpdate(node, m)
	}

	if revision {
		node.UpdatedAt = saved.UpdatedAt
		id := node.Id
		_, err = m.insertNode(node, m.Prefix+"_nodes_audit")

		node.Id = id
		node.UpdatedAt = updatedAt

		helper.PanicOnError(err)
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
		errors.AddError("name", "Name cannot be empty")
	}

	if node.Slug == "" {
		errors.AddError("slug", "Slug cannot be empty")
	}

	if node.Type == "" {
		errors.AddError("type", "Type cannot be empty")
	}

	if len(node.Access) == 0 {
		errors.AddError("access", "Access cannot be empty")
	}

	if node.Status < 0 || node.Status > 3 {
		errors.AddError("status", "Invalid status")
	}

	if h, ok := m.Handlers.Get(node).(ValidateNodeHandler); ok {
		h.Validate(node, m, errors)
	}

	return !errors.HasErrors(), errors
}
