package mock

import (
	"github.com/stretchr/testify/mock"
	sq "github.com/lann/squirrel"
	"container/list"
	nc "github.com/rande/gonode/core"
)

type MockedManager struct {
  mock.Mock
}

func (m *MockedManager) FindBy(query sq.SelectBuilder, offset uint64, limit uint64) *list.List {
	args := m.Mock.Called(query, offset, limit)

	return args.Get(0).(*list.List)
}

func (m *MockedManager) FindOneBy(query sq.SelectBuilder) *nc.Node {
	args := m.Mock.Called(query)

	return args.Get(0).(*nc.Node)
}

func (m *MockedManager) Find(uuid nc.Reference) *nc.Node {
	args := m.Mock.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*nc.Node)
}

func (m *MockedManager) Remove(query sq.SelectBuilder) error {
	args := m.Mock.Called(query)

	return args.Error(0)
}

func (m *MockedManager) RemoveOne(node *nc.Node) (*nc.Node, error) {
	args := m.Mock.Called(node)

	return args.Get(0).(*nc.Node), args.Error(1)
}

func (m *MockedManager) Save(node *nc.Node) (*nc.Node, error) {
	args := m.Mock.Called(node)

	return args.Get(0).(*nc.Node), args.Error(1)
}

func (m *MockedManager) Notify(channel string, payload string) {
	m.Mock.Called(channel, payload)
}

func (m *MockedManager) GetHandler(node *nc.Node) nc.Handler {
	args := m.Mock.Called(node)

	return args.Get(0).(nc.Handler)
}

func (m *MockedManager) NewNode(t string) *nc.Node {
	args := m.Mock.Called(t)

	return args.Get(0).(*nc.Node)
}
