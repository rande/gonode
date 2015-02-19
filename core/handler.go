package core

import (
	pq "github.com/lib/pq"
)

type NodeData interface {}
type NodeMeta interface {}

type Handler interface {
	GetStruct() (NodeData, NodeMeta) // Data, Meta

	PreUpdate(node *Node, m NodeManager) error
	PostUpdate(node *Node, m NodeManager) error
	PreInsert(node *Node, m NodeManager) error
	PostInsert(node *Node, m NodeManager) error
	Validate(node *Node, m NodeManager, e Errors)
}

type Listener interface {
	Handle(notification *pq.Notification, manager NodeManager)
}

const (
	ProcessStatusInit   = 0
	ProcessStatusUpdate = 1
	ProcessStatusDone   = 2
	ProcessStatusError  = 3
)
