BinData
=======

This plugin provides an helper to expose assets to http request.

Configuration
-------------

```toml
[bindata]
    base_path = "/var/go" # default is os.Getenv("GOPATH") + "/src"
    [bindata.assets]
        [bindata.assets.explorer]
        index = "index.html"
        public = "/explorer"
        private = "github.com/rande/gonode/explorer/dist"
```

It is possible to expose any files from the assets packages as describe in the [assets documentations](../../assets/README.md)