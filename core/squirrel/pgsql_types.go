// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package squirrel

// from https://gist.github.com/adharris/4163702

import (
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"regexp"
	"strings"
)

var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	str := string(asBytes)

	if str == "{}" {
		return nil
	}

	// change quote escapes for csv parser
	str = quoteEscapeRegex.ReplaceAllString(str, `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)
	// remove braces
	str = str[1 : len(str)-1]
	csvReader := csv.NewReader(strings.NewReader(str))

	slice, err := csvReader.Read()

	if err != nil {
		return err
	}

	(*s) = StringSlice(slice)

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	// string escapes.
	// \ => \\\
	// " => \"
	for i, elem := range s {
		s[i] = `"` + strings.Replace(strings.Replace(elem, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
