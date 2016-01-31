// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package squirrel

func NewExprSlice(sql string, args ...interface{}) *ExprSlice {
	return &ExprSlice{
		sql:  sql,
		args: args,
	}
}

type ExprSlice struct {
	sql  string
	args []interface{}
}

func (e *ExprSlice) ToSql() (string, []interface{}, error) {

	b := make([]interface{}, 0)

	switch t := e.args[0].(type) {
	case []string:
		for i := range t {
			b = append(b, t[i])
		}
	case []int:
		for i := range t {
			b = append(b, t[i])
		}
	}

	return e.sql, b, nil
}
