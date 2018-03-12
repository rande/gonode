// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockedDriver struct {
	mock.Mock
}

func (m *MockedDriver) Has(key string) bool {
	args := m.Mock.Called(key)

	return args.Bool(0)
}

func (m *MockedDriver) GetReader(key string) (io.ReadCloser, error) {
	args := m.Mock.Called(key)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockedDriver) GetWriter(key string) (io.WriteCloser, error) {
	args := m.Mock.Called(key)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (m *MockedDriver) Remove(key string) error {
	args := m.Mock.Called(key)

	return args.Error(0)
}
