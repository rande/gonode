// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package localisation

import (
	tpl "html/template"
	"testing"
	"time"

	"github.com/rande/gonode/modules/template"
	"github.com/stretchr/testify/assert"
)

func GetFunc(funcs map[string]interface{}, name string) func(ctx template.Context, date time.Time) tpl.HTML {
	return funcs[name].(func(ctx template.Context, date time.Time) tpl.HTML)
}

var funcs = CreateTemplateFuncMap("en_GB")
var date = time.Date(2023, time.December, 25, 1, 1, 1, 1, time.UTC)
var ctx = template.Context{
	"locale": "fr_FR",
	"tz":     "Europe/Paris",
}
var emptyCtx = template.Context{}

func Test_Format_Short_Date(t *testing.T) {
	assert.Equal(t, tpl.HTML("25/12/2023"), GetFunc(funcs, "short_date")(ctx, date))
}

func Test_Format_Short_Date_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("12/25/23"), GetFunc(funcs, "short_date")(emptyCtx, date))
}

func Test_Format_Medium_Date(t *testing.T) {
	assert.Equal(t, tpl.HTML("25 Dec 2023"), GetFunc(funcs, "medium_date")(ctx, date))
}

func Test_Format_Medium_Date_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("Dec 25, 2023"), GetFunc(funcs, "medium_date")(emptyCtx, date))
}

func Test_Format_Long_Date(t *testing.T) {
	assert.Equal(t, tpl.HTML("25 December 2023"), GetFunc(funcs, "long_date")(ctx, date))
}

func Test_Format_Long_Date_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("December 25, 2023"), GetFunc(funcs, "long_date")(emptyCtx, date))
}

func Test_Format_Time(t *testing.T) {
	assert.Equal(t, tpl.HTML("01:01"), GetFunc(funcs, "time")(ctx, date))
}

func Test_Format_Time_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("1:01 AM"), GetFunc(funcs, "time")(emptyCtx, date))
}

func Test_Format_Tz_Time(t *testing.T) {
	assert.Equal(t, tpl.HTML("01:01 (UTC)"), GetFunc(funcs, "tz_time")(ctx, date))
}

func Test_Format_Tz_Time_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("1:01 AM (UTC)"), GetFunc(funcs, "tz_time")(emptyCtx, date))
}

func Test_Format_Tz(t *testing.T) {
	assert.Equal(t, tpl.HTML("Europe/Paris"), GetFunc(funcs, "tz")(ctx, date))
}

func Test_Format_Tz_EmptyCtx(t *testing.T) {
	assert.Equal(t, tpl.HTML("UTC"), GetFunc(funcs, "tz")(emptyCtx, date))
}
