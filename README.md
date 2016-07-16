Go Node
=======

[![Build Status](https://travis-ci.org/rande/gonode.svg?branch=master)](https://travis-ci.org/rande/gonode)
[![Coverage Status](https://coveralls.io/repos/github/rande/gonode/badge.svg)](https://coveralls.io/github/rande/gonode)

A prototype to store dynamic node inside a PostgreSQL database with the JSONb storage column.

Documentation
-------------
 
 * [Install](docs/install.md)
 * [Contributing](docs/contributing.md)
 * [Core](docs/core)
    * [Router](docs/core/router.md): named routes with goji
    * [Vault](docs/core/vault.md): Binary storage with secure option
    * [Guard](docs/core/guard.md): Authentification
    * [Security](docs/core/security.md): CORS 
    * [Logger](docs/core/logger.md): Logger
    * [Bindata](docs/core/bindata.md): Provide http handler to server static file from ``go-bindata`` assets
 * [Modules](docs/modules)
    * [Node](docs/modules/node.md): Node principles
    * [Search](docs/modules/search.md): Search filters
    * [Restful API](docs/modules/restful_api.md): Restful API 
    * [Raw](docs/modules/raw.md): Send raw content
    * [Prism](docs/modules/prism.md)
    * [Feed](docs/modules/feed.md)
    * [Media](docs/modules/media.md)
    * [Access](docs/modules/access.md)
