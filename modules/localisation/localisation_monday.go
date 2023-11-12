// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package localisation

import (
	tpl "html/template"
	"time"

	"github.com/goodsign/monday"
	"github.com/rande/gonode/modules/template"
)

func GetLocale(ctx template.Context) monday.Locale {
	locale := "en_US"
	if l, ok := ctx["locale"]; ok {
		locale = l.(string)
	}

	return monday.Locale(locale)
}

func GetTimezone(ctx template.Context) *time.Location {
	tz := "UTC"
	if l, ok := ctx["tz"]; ok {
		tz = l.(string)
	}

	if location, err := time.LoadLocation(tz); err == nil {
		return location
	}

	return time.UTC
}

func CreateTemplateFuncMap(defaultLocale string) map[string]interface{} {

	FuncMap := make(map[string]interface{})

	FuncMap["short_date"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(date.Format(monday.ShortFormatsByLocale[GetLocale(ctx)]))
	}

	FuncMap["medium_date"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(date.Format(monday.MediumFormatsByLocale[GetLocale(ctx)]))
	}

	FuncMap["long_date"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(date.Format(monday.LongFormatsByLocale[GetLocale(ctx)]))
	}

	FuncMap["time"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(date.In(GetTimezone(ctx)).Format(monday.TimeFormatsByLocale[GetLocale(ctx)]))
	}

	FuncMap["tz_time"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(date.In(GetTimezone(ctx)).Format(monday.TimeFormatsByLocale[GetLocale(ctx)] + " (MST)"))
	}

	FuncMap["tz"] = func(ctx template.Context, date time.Time) tpl.HTML {
		return tpl.HTML(GetTimezone(ctx).String())
	}

	return FuncMap
}
