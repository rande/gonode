package gonode

import (
	"time"
	"github.com/twinj/uuid"
)

type Node struct {
	id         int
	Type       string
	Name       string
	Slug       string
	Data       interface {}
	Meta       interface {}
	Status     int
	Weight     int
	Revision   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Enabled    bool
	Deleted    bool
	Set        string
	Parents    []string
	Uuid       Uuid
	UpdatedBy  Uuid
	CreatedBy  Uuid
	ParentUuid Uuid
	SetUuid    Uuid
	Source     Uuid
}

func (node *Node) Id() int {
	return node.id
}

func NewNode() *Node {
	return &Node{
		Uuid:       GetEmptyUuid(),
		Source:     GetEmptyUuid(),
		ParentUuid: GetEmptyUuid(),
		UpdatedBy:  GetEmptyUuid(),
		CreatedBy:  GetEmptyUuid(),
		SetUuid:    GetEmptyUuid(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Weight:     1,
		Revision:   1,
		Deleted:    false,
		Enabled:    true,
		Status:     StatusDraft,
	}
}
