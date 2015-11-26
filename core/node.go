// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"fmt"
	"github.com/twinj/uuid"
	"time"
)

var (
	StatusNew       = 0
	StatusDraft     = 1
	StatusCompleted = 2
	StatusValidated = 3
)

type Reference struct {
	uuid.UUID
}

func (m *Reference) MarshalJSON() ([]byte, error) {
	// Manually calling Marshal for Contents
	cont, err := json.Marshal(uuid.Formatter(m.UUID, uuid.CleanHyphen))
	if err != nil {
		return nil, err
	}

	// Stitching it all together
	return cont, nil
}

func (m *Reference) UnmarshalJSON(data []byte) error {
	PanicIf(len(data) < 32, "invalid uuid size")

	tmpUuid, err := uuid.Parse(string(data[1 : len(data)-1]))

	if err != nil {
		return err
	}

	m.UUID = GetReference(tmpUuid)

	return nil
}

func (m *Reference) CleanString() string {
	return uuid.Formatter(m.UUID, uuid.CleanHyphen)
}

func GetReferenceFromString(reference string) (Reference, error) {
	v, err := uuid.Parse(reference)

	if err != nil {
		return GetEmptyReference(), InvalidReferenceFormatError
	}

	return GetReference(v), nil
}

func GetReference(uuid uuid.UUID) Reference {
	return Reference{uuid}
}

type Node struct {
	id         int
	Uuid       Reference   `json:"uuid"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta"`
	Status     int         `json:"status"`
	Weight     int         `json:"weight"`
	Revision   int         `json:"revision"`
	Version    int         `json:"version"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Enabled    bool        `json:"enabled"`
	Deleted    bool        `json:"deleted"`
	Parents    []Reference `json:"parents"`
	UpdatedBy  Reference   `json:"updated_by"`
	CreatedBy  Reference   `json:"created_by"`
	ParentUuid Reference   `json:"parent_uuid"`
	SetUuid    Reference   `json:"set_uuid"`
	Source     Reference   `json:"source"`
}

func (node *Node) Id() int {
	return node.id
}

func (node *Node) UniqueId() string {
	return fmt.Sprintf("%s-v%d", node.Uuid.CleanString(), node.Revision)
}

func NewNode() *Node {
	return &Node{
		Uuid:       GetEmptyReference(),
		Source:     GetEmptyReference(),
		ParentUuid: GetEmptyReference(),
		UpdatedBy:  GetEmptyReference(),
		CreatedBy:  GetEmptyReference(),
		SetUuid:    GetEmptyReference(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Weight:     1,
		Revision:   1,
		Version:    1,
		Deleted:    false,
		Enabled:    true,
		Status:     StatusNew,
	}
}

func DumpNode(node *Node) {
	PanicIf(node == nil, "Cannot dump, node is nil")

	fmt.Printf(" >>> Node: %+v\n", node.id)
	fmt.Printf(" Uuid:       %s\n", node.Uuid)
	fmt.Printf(" Type:       %s\n", node.Type)
	fmt.Printf(" Name:       %s\n", node.Name)
	fmt.Printf(" Status:     %d\n", node.Status)
	fmt.Printf(" Weight:     %d\n", node.Weight)
	fmt.Printf(" Deleted:    %t\n", node.Deleted)
	fmt.Printf(" Enabled:    %t\n", node.Enabled)
	fmt.Printf(" Revision:   %d\n", node.Revision)
	fmt.Printf(" Version:    %d\n", node.Version)
	fmt.Printf(" CreatedAt:  %+v\n", node.CreatedAt)
	fmt.Printf(" UpdatedAt:  %+v\n", node.UpdatedAt)
	fmt.Printf(" Slug:       %s\n", node.Slug)
	fmt.Printf(" Data:       %T => %+v\n", node.Data, node.Data)
	fmt.Printf(" Meta:       %T => %+v\n", node.Meta, node.Meta)
	fmt.Printf(" CreatedBy:  %s\n", node.CreatedBy)
	fmt.Printf(" UpdatedBy:  %s\n", node.UpdatedBy)
	fmt.Printf(" ParentUuid: %s\n", node.ParentUuid)
	fmt.Printf(" Parents:    %+v\n", node.Parents)
	fmt.Printf(" SetUuid:    %s\n", node.SetUuid)
	fmt.Printf(" Source:     %s\n", node.Source)
	fmt.Printf(" <<< End Node\n")
}
