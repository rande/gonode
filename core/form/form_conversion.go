// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"fmt"
	"reflect"
	"strconv"
)

func StrToValue(value interface{}, src reflect.Value) (interface{}, bool) {
	if value == nil {
		return nil, false
	}

	switch src.Kind() {
	case reflect.Bool:
		return StrToBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return StrToNumber(value.(string), src.Kind())
	case reflect.String:
		return value, true
	default:
		panic(fmt.Sprintf("Case not implemented - StrToValue: value `%s` to type `%s`", value, src.Kind()))
	}
}

func ValueToStrSlice(value interface{}, src reflect.Value) ([]string, bool) {
	c := make([]string, 0)
	// for _, _ := range value.([]interface{}) {
	// 	// c = append(c, ValueToStr(v, kind.Elem().Kind().Elem()))
	// }
	return c, true
}

func ValueToStr(value interface{}, src reflect.Value) (string, bool) {
	if value == nil {
		return "", false
	}

	//fmt.Printf("ValueToStr: %s - %s, %s\n", value, src.Kind(), src.Kind())

	switch src.Kind() {

	case reflect.Bool:
		return BoolToStr(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return NumberToStr(value)
	case reflect.String:
		return value.(string), true
	default:
		panic(fmt.Sprintf("Case not implemented - ValueToStr: value `%s` to type `%s`", value, src.Kind()))
	}
}

func yes(value string) bool {
	return value == "checked" || value == "true" || value == "1" || value == "on" || value == "yes"
}

func no(value string) bool {
	return value == "false" || value == "0" || value == "off" || value == "no"
}

func BoolToStr(v interface{}) (string, bool) {
	if v, ok := v.(bool); ok {
		if v {
			return "true", true
		} else {
			return "false", true
		}
	}

	return "", false
}

func StrToBool(value interface{}) (bool, bool) {
	if value == nil {
		return false, false
	}

	if yes(value.(string)) {
		return true, true
	} else if no(value.(string)) {
		return false, true
	}

	return false, false
}

func NumberToStr(v interface{}) (string, bool) {
	return fmt.Sprintf("%v", v), true
}

func StrToNumber(value string, kind reflect.Kind) (interface{}, bool) {
	if kind == reflect.Int {
		if i, err := strconv.ParseInt(value, 10, 0); err == nil {
			return int(i), true
		}
	} else if kind == reflect.Int8 {
		if i, err := strconv.ParseInt(value, 10, 8); err == nil {
			return int8(i), true
		}
	} else if kind == reflect.Int16 {
		if i, err := strconv.ParseInt(value, 10, 16); err == nil {
			return int16(i), true
		}
	} else if kind == reflect.Int32 {
		if i, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(i), true

		}
	} else if kind == reflect.Int64 {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64(i), true

		}
	} else if kind == reflect.Uint {
		if i, err := strconv.ParseUint(value, 10, 0); err == nil {
			return uint(i), true

		}
	} else if kind == reflect.Uint8 {
		if i, err := strconv.ParseUint(value, 10, 8); err == nil {
			return uint8(i), true

		}
	} else if kind == reflect.Uint16 {
		if i, err := strconv.ParseUint(value, 10, 16); err == nil {
			return uint16(i), true

		}
	} else if kind == reflect.Uint32 {
		if i, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(i), true

		}
	} else if kind == reflect.Uint64 {
		if i, err := strconv.ParseUint(value, 10, 64); err == nil {
			return uint64(i), true

		}
	} else if kind == reflect.Float32 {
		if i, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(i), true

		}
	} else if kind == reflect.Float64 {
		if i, err := strconv.ParseFloat(value, 64); err == nil {
			return float64(i), true

		}
	}

	return nil, false
}
