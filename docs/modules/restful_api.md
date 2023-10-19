# Restful API

The API provides information about the current system, but also a way to interact with nodes.

## Version

The API version is part of the requested URL: `/api/:version/...`, for now there is only on version `v1.0`. Any
other version will generate a `Bad Request`.

## Node API

-   Create one node
    -   method: `POST /api/:version/nodes`
    -   role: `node:api:create`
-   Get one node
    -   method: `GET /api/:version/nodes/:nid`
    -   role: `node:api:read`
-   Alter one node
    -   method: `PUT /api/:version/nodes/:nid`
    -   role: `node:api:update`
-   Delete one node
    -   method: `DELETE /api/:version/nodes/:nid`
    -   role: `node:api:delete`
-   Get node revisions
    -   method: `GET /api/:version/nodes/:nid/revisions`
    -   role: `node:api:revisions`
-   Get one node revision
    -   method: `GET /api/:version/nodes/:nid/revisions/:rev`
    -   role: `node:api:revision`
-   Move `nid` as a child of `parentNid`
    -   method: `PUT /api/:version/nodes/move/:nid/:parentNid`
    -   role: `node:api:move`
-   List node (see [search.md](search.md))
    -   method: `GET /api/:version/nodes`
    -   role: `node:api:list`
-   Basic url to return hello
    -   method: `GET /api/:version/hello`
    -   role: `-`
-   Notify subscribers
    -   method: `PUT /api/:version/notify/:name`
    -   role: `node:api:notify`
-   Websocket to retrieve update stream
    -   method: `WS /api/:version/nodes/stream`
    -   role: `node:api:stream`

Please note: the `node:api:master` role will allow any actions to be performed.

## Instrospection API

-   `GET /:version/handlers/node`: return a list of node handlers
-   `GET /:version/handlers/view`: return a list of view node handlers
-   `GET /:version/services`: return a list of services

## Security

The security is based on the [guard module](../core/guard.md). Please see related documentation.
