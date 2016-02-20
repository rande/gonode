Logger
======

Integrate logrus into gonode.

Configuration
-------------

```toml
[logger]
    level = "debug"

    [logger.fields]
    app = "gonode"

    [logger.hooks]
        [logger.hooks.default]
        service = "influxdb"
        dsn = "http://localhost:8086"
        tags = ["app.core"]
        database = "stats"
        level = "debug"
```

Feel free to send PR to add support for others hooks.

Usage
-----

Default usage:

```golang
import (
    log "github.com/Sirupsen/logrus"
)

logger := app.Get("logger").(*log.Logger)
logger.WithFields(log.Fields{
    "type":   node.Type,
    "uuid":   node.Uuid,
    "module": "core.manager",
}).Warn("soft delete one")

```
  
Request's logger:

```golang

mux.Get(publicPath+"/*", func(c web.C, res http.ResponseWriter, req *http.Request) {
    var logger *log.Entry

    if l, ok := c.Env["logger"]; ok {
        logger = l.(*log.Entry).WithFields(log.Fields{
            "module": "bindata.handler",
        })
    }
})

```