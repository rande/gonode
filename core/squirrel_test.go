// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	sq "github.com/lann/squirrel"
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/twinj/uuid"
)

func Test_ExprSlice_With_int(t *testing.T) {
	value := []int{1, 2, 3}

	e := NewExprSlice("data->'%s' ??| array["+sq.Placeholders(len(value))+"]", value)

	sql, args, error := e.ToSql()

	assert.Nil(t, error)
	assert.Equal(t, "data->'%s' ??| array[?,?,?]", sql)
	assert.Equal(t, []interface{}{1, 2, 3}, args)
}

func Test_ExprSlice_With_string(t *testing.T) {

	value := []string{"1", "2", "3"}

	e := NewExprSlice("data->'%s' ??| array["+sq.Placeholders(len(value))+"]", value)

	sql, args, error := e.ToSql()

	assert.Nil(t, error)
	assert.Equal(t, "data->'%s' ??| array[?,?,?]", sql)
	assert.Equal(t, []interface{}{"1", "2", "3"}, args)
}
