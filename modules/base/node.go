// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"errors"
	"fmt"
	"time"

	"github.com/rande/gonode/core/helper"
)

var (
	StatusNew       = 0
	StatusDraft     = 1
	StatusCompleted = 2
	StatusValidated = 3
)

func GetReference(nid string) string {
	return nid
}

type Modules map[string]interface{}

func (p Modules) Set(name string, v interface{}) {
	p[name] = v
}

func (p Modules) Has(name string) bool {
	if _, ok := p[name]; ok {
		return true
	}

	return false
}

func (p Modules) Get(name string) (interface{}, error) {
	if p.Has(name) {
		return p[name], nil
	}

	return nil, errors.New("No modules")
}

type Node struct {
	Id        int         `json:"-"`
	Nid       string      `json:"nid"`
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	Path      string      `json:"path"`
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta"`
	Status    int         `json:"status"`
	Weight    int         `json:"weight"`
	Revision  int         `json:"revision"`
	Version   int         `json:"version"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Enabled   bool        `json:"enabled"`
	Deleted   bool        `json:"deleted"`
	Parents   []string    `json:"parents"`
	UpdatedBy string      `json:"updated_by"`
	CreatedBy string      `json:"created_by"`
	ParentNid string      `json:"parent_nid"`
	SetNid    string      `json:"set_nid"`
	Source    string      `json:"source"`
	Modules   Modules     `json:"modules"`
	Access    []string    `json:"access"` // key => roles required to access the nodes
}

func (node *Node) UniqueId() string {
	return fmt.Sprintf("%s-v%d", node.Nid, node.Revision)
}

func NewNode() *Node {
	return &Node{
		Nid:       GetEmptyReference(),
		Source:    GetEmptyReference(),
		ParentNid: GetEmptyReference(),
		UpdatedBy: GetEmptyReference(),
		CreatedBy: GetEmptyReference(),
		SetNid:    GetEmptyReference(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Weight:    1,
		Revision:  1,
		Version:   1,
		Deleted:   false,
		Enabled:   true,
		Status:    StatusNew,
		Modules:   make(map[string]interface{}, 0),
		Access:    make([]string, 0),
	}
}

func DumpNode(node *Node) {
	helper.PanicIf(node == nil, "Cannot dump, node is nil")

	fmt.Printf(" >>> Node: %v\n", node.Id)
	fmt.Printf(" Nid:       %s\n", node.Nid)
	fmt.Printf(" Type:       %s\n", node.Type)
	fmt.Printf(" Name:       %s\n", node.Name)
	fmt.Printf(" Status:     %d\n", node.Status)
	fmt.Printf(" Weight:     %d\n", node.Weight)
	fmt.Printf(" Deleted:    %t\n", node.Deleted)
	fmt.Printf(" Enabled:    %t\n", node.Enabled)
	fmt.Printf(" Revision:   %d\n", node.Revision)
	fmt.Printf(" Version:    %d\n", node.Version)
	fmt.Printf(" CreatedAt:  %v\n", node.CreatedAt)
	fmt.Printf(" UpdatedAt:  %v\n", node.UpdatedAt)
	fmt.Printf(" Slug:       %s\n", node.Slug)
	fmt.Printf(" Path:       %s\n", node.Path)
	fmt.Printf(" Data:       %T => %v\n", node.Data, node.Data)
	fmt.Printf(" Meta:       %T => %v\n", node.Meta, node.Meta)
	fmt.Printf(" Modules:    %T => %v\n", node.Modules, node.Modules)
	fmt.Printf(" CreatedBy:  %s\n", node.CreatedBy)
	fmt.Printf(" UpdatedBy:  %s\n", node.UpdatedBy)
	fmt.Printf(" ParentNid: %s\n", node.ParentNid)
	fmt.Printf(" Parents:    %v\n", node.Parents)
	fmt.Printf(" SetNid:    %s\n", node.SetNid)
	fmt.Printf(" Source:     %s\n", node.Source)
	fmt.Printf(" <<< End Node\n")
}
