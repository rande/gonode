// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_Dispatch_SameLevel(t *testing.T) {

	e := &log.Entry{}
	e.Level = log.DebugLevel

	h := &MockedHook{}
	h.On("Fire", e).Return(nil)

	d := &DispatchHook{
		Hooks: make(map[log.Level][]log.Hook, 0),
	}
	d.Add(h, log.DebugLevel)
	d.Fire(e)

	h.AssertCalled(t, "Fire", e)
}

func Test_Dispatch_Debug(t *testing.T) {

	e := &log.Entry{}
	e.Level = log.DebugLevel

	h := &MockedHook{}
	h.On("Fire", e).Return(nil)

	d := &DispatchHook{
		Hooks: make(map[log.Level][]log.Hook, 0),
	}
	d.Add(h, log.WarnLevel)
	d.Fire(e)

	h.AssertNotCalled(t, "Fire", e)
}

func Test_Dispatch_Fatal(t *testing.T) {

	e := &log.Entry{}
	e.Level = log.FatalLevel

	h := &MockedHook{}
	h.On("Fire", e).Return(nil)

	d := &DispatchHook{
		Hooks: make(map[log.Level][]log.Hook, 0),
	}
	d.Add(h, log.WarnLevel)
	d.Fire(e)

	h.AssertCalled(t, "Fire", e)
}
