package extra

import (
	sq "github.com/rande/squirrel"
)

func ExprSlice(sql string, size int, args ...interface{}) sq.Expression {
	if len(args) != 1 {
		return sq.Expr(sql, args)
	}

	b := make([]interface{}, size)

	switch t := args[0].(type) {
	case []string:
		for i := range t {
			b[i] = t[i]
		}
	case []int:
		for i := range t {
			b[i] = t[i]
		}
	}

	return sq.Expr(sql, b)
}
