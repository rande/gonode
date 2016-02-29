Restful API
===========

The API provides information about the current system, but also a way to interact with nodes.

## Version

The API version is part of the requested URL: ``/api/:version/...``, for now there is only on version ``v1.0``. Any
other version will generate a ``Bad Request``.

## Node API 

 - ``POST /api/:version/nodes``: create one node 
 - ``GET /api/:version/nodes/:uuid``: get one node
 - ``PUT /api/:version/nodes/:uuid``: alter one node
 - ``DELETE /api/:version/nodes/:uuid``: delete one node
 - ``GET /api/:version/nodes/:uuid/revisions``: get node revisions
 - ``GET /api/:version/nodes/:uuid/revisions/:rev``: get one node revision
 - ``PUT /api/:version/nodes/move/:uuid/:parentUuid``: move ``uuid`` as a child of ``parentUuid`` 
 - ``GET /api/:version/nodes``: list node (see [search.md](search.md)) 
 - ``GET /api/:version/hello``: basic url to return hello
 - ``PUT /api/:version/notify/:name``: notify subscribers
 - ``WS /api/:version/nodes/stream``: websocket to retrieve update stream
 
## Instrospection API

 - ``GET /:version/handlers/node``: return a list of node handlers 
 - ``GET /:version/handlers/view``: return a list of view node handlers
 - ``GET /:version/services``: return a list of services


## Security

The security is based on the [guard module](../core/guard.md). Please see related documentation.
