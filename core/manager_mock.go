// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"container/list"
	sq "github.com/lann/squirrel"
	"github.com/stretchr/testify/mock"
)

type MockedManager struct {
	mock.Mock
}

func (m *MockedManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	args := m.Mock.Called(query, offset, limit)

	return args.Get(0).(*list.List)
}

func (m *MockedManager) FindOneBy(query sq.SelectBuilder) *Node {
	args := m.Mock.Called(query)

	return args.Get(0).(*Node)
}

func (m *MockedManager) Find(uuid Reference) *Node {
	args := m.Mock.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*Node)
}

func (m *MockedManager) Remove(query sq.SelectBuilder) error {
	args := m.Mock.Called(query)

	return args.Error(0)
}

func (m *MockedManager) RemoveOne(node *Node) (*Node, error) {
	args := m.Mock.Called(node)

	return args.Get(0).(*Node), args.Error(1)
}

func (m *MockedManager) Save(node *Node, revision bool) (*Node, error) {
	args := m.Mock.Called(node)

	return args.Get(0).(*Node), args.Error(1)
}

func (m *MockedManager) Notify(channel string, payload string) {
	m.Mock.Called(channel, payload)
}

func (m *MockedManager) NewNode(t string) *Node {
	args := m.Mock.Called(t)

	return args.Get(0).(*Node)
}

func (m *MockedManager) SelectBuilder() sq.SelectBuilder {
	args := m.Mock.Called()

	return args.Get(0).(sq.SelectBuilder)
}

func (m *MockedManager) Validate(node *Node) (bool, Errors) {
	args := m.Mock.Called(node)

	return args.Get(0).(bool), args.Get(1).(Errors)
}

func (m *MockedManager) Move(uuid, parentUuid Reference) (int64, error) {
	args := m.Mock.Called(uuid, parentUuid)

	return args.Get(0).(int64), args.Error(1)
}
