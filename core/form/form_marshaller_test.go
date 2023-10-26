// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"testing"
	"time"

	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
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

	assert.Equal(t, expectedDate, field.SubmittedValue)
}

func Test_Date_Unmarshalling_With_Time(t *testing.T) {
	node := base.NewNode()
	handler := &blog.PostHandler{}
	node.Data, node.Meta = handler.GetStruct()

	form := CreateForm(node)
	form.Add("UpdatedAt", "date")

	dataForm := CreateForm(node.Data)
	dataForm.Add("PublicationDate", "date")

	form.Add("data", "form", dataForm)

	PrepareForm(form)

	v := url.Values{
		"UpdatedAt": []string{"2022-04-01"},
		"data.PublicationDate": []string{
			"2022-04-01",
		},
	}

	BindUrlValues(form, v)

	ValidateForm(form)

	AttachValues(form)

	assert.Equal(t, time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC), node.UpdatedAt)
	assert.Equal(t, time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC), node.Data.(*blog.Post).PublicationDate)
}
