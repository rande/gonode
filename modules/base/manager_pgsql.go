// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/core/squirrel"
	log "github.com/sirupsen/logrus"
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
		SelectClause: "id, nid, type, name, revision, version, created_at, updated_at, set_nid, parent_nid, parents, slug, path, created_by, updated_by, data, meta, modules, access, deleted, enabled, source, status, weight",
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

func (m *PgNodeManager) Find(nid string) *Node {
	return m.FindOneBy(m.SelectBuilder(NewSelectOptions()).Where(sq.Eq{"nid": nid}))
}

func (m *PgNodeManager) hydrate(rows *sql.Rows) *Node {
	node := &Node{}

	data := json.RawMessage{}
	meta := json.RawMessage{}
	modules := json.RawMessage{}

	Nid := ""
	SetNid := ""
	ParentNid := ""
	CreatedBy := ""
	UpdatedBy := ""
	Source := ""

	var Access, Parents squirrel.StringSlice

	err := rows.Scan(
		&node.Id,
		&Nid,
		&node.Type,
		&node.Name,
		&node.Revision,
		&node.Version,
		&node.CreatedAt,
		&node.UpdatedAt,
		&SetNid,
		&ParentNid,
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

	node.Nid = Nid
	node.SetNid = SetNid
	node.ParentNid = ParentNid
	node.CreatedBy = CreatedBy
	node.UpdatedBy = UpdatedBy
	node.Source = Source

	pNids := make([]string, 0)

	for _, pNid := range Parents {
		pNids = append(pNids, GetReference(pNid))
	}

	node.Parents = pNids

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
				Subject:  node.Nid,
				Revision: node.Revision,
				Date:     node.UpdatedAt,
			})

			if m.Logger != nil {
				m.Logger.WithFields(log.Fields{
					"type":   node.Type,
					"nid":    node.Nid,
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
			"nid":    node.Nid,
			"module": "node.manager",
		}).Warn("soft delete one")
	}

	m.sendNotification(m.Prefix+"_manager_action", &ModelEvent{
		Type:     node.Type,
		Action:   "SoftDelete",
		Subject:  node.Nid,
		Revision: node.Revision,
		Date:     node.UpdatedAt,
		Name:     node.Name,
	})

	return m.Save(node, true)
}

func (m *PgNodeManager) insertNode(node *Node, table string) (*Node, error) {
	if node.Nid == GetEmptyReference() {
		node.Nid = NewId()
	}

	if node.Slug == "" {
		node.Slug = node.Nid
	}

	Parents := make(squirrel.StringSlice, 0)
	for _, p := range node.Parents {
		Parents = append(Parents, p)
	}

	node.Access = security.EnsureRoles(node.Access, "node:api:master")

	Access := make(squirrel.StringSlice, 0)
	for _, a := range node.Access {
		Access = append(Access, a)
	}

	query := sq.Insert(table).
		Columns(
			"nid", "type", "revision", "version", "name", "created_at", "updated_at", "set_nid",
			"parent_nid", "parents", "slug", "path", "created_by", "updated_by", "data", "meta", "modules",
			"access", "deleted", "enabled", "source", "status", "weight").
		Values(
			node.Nid,
			node.Type,
			node.Revision,
			node.Version,
			node.Name,
			node.CreatedAt,
			node.UpdatedAt,
			node.SetNid,
			node.ParentNid,
			Parents,
			node.Slug,
			node.Path,
			node.CreatedBy,
			node.UpdatedBy,
			string(InterfaceToJsonMessage(node.Type, node.Data)[:]),
			string(InterfaceToJsonMessage(node.Type, node.Meta)[:]),
			string(InterfaceToJsonMessage(node.Type, node.Modules)[:]),
			Access,
			node.Deleted,
			node.Enabled,
			node.Source,
			node.Status,
			node.Weight,
		).
		Suffix("RETURNING \"id\"").
		RunWith(m.Db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&node.Id)

	return node, err
}

func (m *PgNodeManager) Move(nid, parentNid string) (int64, error) {
	tx, err := m.Db.Begin()

	if err != nil {
		return 0, err
	}

	r, err := tx.Exec(fmt.Sprintf(`UPDATE %s SET parent_nid = $1 WHERE nid = $2 AND EXISTS(SELECT nid FROM %s WHERE nid = $3 and $4 <> ALL(parents))`,
		m.Prefix+"_nodes", m.Prefix+"_nodes"),
		parentNid,
		nid,
		parentNid,
		nid)

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
					SELECT nid, parent_nid, parents,
						CASE	WHEN type = 'core.root' THEN ''
							WHEN array_length(parents, 1)>0 THEN path
							ELSE '/' || slug::varchar(2000)
						END as path
					FROM %s r
					WHERE nid = $1::nid
				UNION ALL
					SELECT c.nid, c.parent_nid, array_append(r.parents, c.parent_nid) AS parents, r.path || '/' || c.slug as path
					FROM %s c
					JOIN r ON c.parent_nid = r.nid
			)
			UPDATE %s n SET parents = r.parents, path = r.path FROM r WHERE r.nid = n.nid`,
			m.Prefix+"_nodes",
			m.Prefix+"_nodes",
			m.Prefix+"_nodes"),
			parentNid)
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
		Set("nid", node.Nid).
		Set("type", node.Type).
		Set("revision", node.Revision).
		Set("version", node.Version).
		Set("name", node.Name).
		Set("created_at", node.CreatedAt).
		Set("updated_at", node.UpdatedAt).
		Set("set_nid", node.SetNid).
		Set("slug", node.Slug).
		Set("path", node.Path).
		Set("created_by", node.CreatedBy).
		Set("updated_by", node.UpdatedBy).
		Set("deleted", node.Deleted).
		Set("enabled", node.Enabled).
		Set("data", string(InterfaceToJsonMessage(node.Type, node.Data)[:])).
		Set("meta", string(InterfaceToJsonMessage(node.Type, node.Meta)[:])).
		Set("modules", string(InterfaceToJsonMessage(node.Type, node.Modules)[:])).
		Set("access", Access).
		Set("source", node.Source).
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

func (m *PgNodeManager) Save(node *Node, newRevision bool) (*Node, error) {

	var contextLogger *log.Entry

	if m.Logger != nil {
		contextLogger = m.Logger.WithFields(log.Fields{
			"nid":      node.Nid,
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
			Subject:     node.Nid,
			Date:        node.CreatedAt,
			Name:        node.Name,
			Revision:    node.Revision,
			NewRevision: newRevision,
		})

		return node, err
	}

	if h, ok := handler.(DatabaseNodeHandler); ok {
		h.PreUpdate(node, m)
	}

	// 1. check if the one in the datastore is older
	saved := m.FindOneBy(m.SelectBuilder(NewSelectOptions()).Where(sq.Eq{"nid": node.Nid}))

	if saved != nil && node.Revision != saved.Revision {
		if contextLogger != nil {
			contextLogger.Info("invalid revision for node")
		}

		return node, NewRevisionError(fmt.Sprintf("Invalid revision for node: %s, saved rev: %d, current rev: %d", node.Nid, saved.Revision, node.Revision))
	}

	if contextLogger != nil {
		contextLogger.Debug("updating node")
	}

	if newRevision {
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

	if newRevision {
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
		Subject:     node.Nid,
		Revision:    node.Revision,
		Date:        node.UpdatedAt,
		Name:        node.Name,
		NewRevision: newRevision,
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
