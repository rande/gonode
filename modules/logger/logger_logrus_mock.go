// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedHook struct {
	mock.Mock
}

func (m *MockedHook) Levels() []log.Level {

	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	}
}

func (m *MockedHook) Fire(e *log.Entry) error {
	args := m.Mock.Called(e)

	return args.Error(0)
}
