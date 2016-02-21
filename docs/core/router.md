Router
======

This modules allows to created named routes for goji handlers. There are 3 types of routes:

- ``ABSOLUTE_URL``: this will generate a string like ``http://myserver.com/hello/world``  
- ``ABSOLUTE_PATH``: this will generate a string like ``/hello/world``
- ``NETWORK_PATH``: this will generate a string like ``//hello/world``

The ``ABSOLUTE_URL`` option generates links depends on the current request's information available as a RequestContext object.
The request context provides resolved information about how to generate a valid url depends on reverse proxy information.
The context is created by a middleware.

Usage
-----

Template's helpers:

```jinja
{{ path("prism", url_values("uuid", elm.Uuid )) }} => "/prism/21779d51-122c-4ea9-a09e-9685610adc5c"
{{ url("prism", url_values("uuid", elm.Uuid ), request_context) }} => "http://localhost:2405/prism/21779d51-122c-4ea9-a09e-9685610adc5c"
{{ net("prism", url_values("uuid", elm.Uuid )) }}  => "//prism/21779d51-122c-4ea9-a09e-9685610adc5c"
```

Handler usage:

```go
package example

import (
	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"net/http"
	"github.com/rande/gonode/modules/router"
)

func MyHandler(c web.C, res http.ResponseWriter, req *http.Request) {
    context = &pongo2.Context{}
    
    if _, ok :=  c.Env["request_context"]; ok {
        context["request_context"] = c.Env["request_context"]
    } else {
        context["request_context"] = nil
    }

    tpl, _ := pongo.FromFile("mytemplate.tpl")

    data, _ := tpl.ExecuteWriter(context, res)
}

```

And the related template:

```jinja
<?xml version="1.0" ?>
<rss version="2.0">
    <channel>
        <title>{{ node.Data.Title }}</title>
        <link>{{ request.URL }}</link>
        <description>{{ node.Data.description }}</description>

        {% for elm in pager.Elements %}
            <item>
                 <title>{{ node_data(node, "name") }}</title>
                 <link>{{ url("prism", url_values("uuid", elm.Uuid), request_context) }}</link>
                 <description><![CDATA[{{ node_data(node, "description") }}]]></description>
                 <pubDate>{{ node_data(node, "publication_date") }}</pubDate>
                 <gui>{{ path("prism", url_values("uuid", elm.Uuid)) }}</gui>
            </item>
        {% endfor %}
    </channel>
</rss>
```

Finally how to register the route:

```go

package prism

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/router"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		r := app.Get("gonode.router").(*router.Router)

		r.Handle("prism", "/prism/:uuid", MyHandler)

		return nil
	})
}
```