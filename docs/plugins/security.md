Security
========

CORS
----

The project integrated [https://github.com/rs/cors](https://github.com/rs/cors) to handle CORS security options. For more
information about CORS, please review [Cross Origin Resource Sharing W3 specification](http://www.w3.org/TR/cors/).

### Configuration

```toml
[security]
    [security.cors]
    allowed_origins = ["*"]
    allowed_methods = ["GET", "PUT", "POST"]
    allowed_headers = ["Origin", "Accept", "Content-Type", "Authorization"]
    exposes_headers = []
    allow_credentials = false
    max_age = 0
    options_passthrough = false
```

- ``allowed_origins``: A list of origins a cross-domain request can be executed from. If the special * value is present in the list, all origins will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality. Only one wildcard can be used per origin. The default value is *.
- ``allowed_methods``: A list of methods the client is allowed to use with cross-domain requests. Default value is simple methods (GET and POST).
- ``allowed_headers``: A list of non simple headers the client is allowed to use with cross-domain requests.
- ``exposes_headers``: Indicates which headers are safe to expose to the API of a CORS API specification
- ``allow_credentials``: Indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates. The default is false.
- ``max_age``: Indicates how long (in seconds) the results of a preflight request can be cached. The default is 0 which stands for no max age.
- ``options_passthrough``: Instructs preflight to let other potential next handlers to process the OPTIONS method. Turn this on if your application handles OPTIONS.