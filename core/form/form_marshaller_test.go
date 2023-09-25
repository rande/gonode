// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Date_Marshalling(t *testing.T) {
	field := CreateFormField()
	form := CreateForm(nil)

	field.InitialValue = time.Date(2022, time.April, 1, 1, 1, 1, 1, time.UTC)

	dateMarshal(field, form)

	assert.Equal(t, "2022-04-01", field.Input.Value)
}

func Test_Date_Unmarshalling(t *testing.T) {
	field := CreateFormField()
	field.Input.Id = "Date"
	field.Input.Name = "Date"

	form := CreateForm(nil)

	field.InitialValue = time.Date(2022, time.April, 1, 1, 1, 1, 1, time.UTC)

	v := url.Values{
		"Date": []string{"2022-04-01"},
	}

	dateUnmarshal(field, form, v)

	expectedDate := time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, expectedDate, field.SubmitedValue)
}