# Prism

This plugin allows to render a node calling a specific view handler for the current node type.

## Request / Response Workflow

When the prism controller is being called a `base.ViewRequest` is created holding the current `http.Request`
structure, also a `base.ViewResponse` with the `ResponseWriter`. So the ViewHandler is called with:

`Execute(node *Node, request *ViewRequest, response *ViewResponse) error`

The `Execute` method can either use the `ResponseWriter` to write content directly to the client or set the
`Template` and the `Context` property from the `ViewResponse` structure.

If the template is set, the controller will use this template to generate the related content. The `ViewHandler`
is like a small controller dedicated to one node.

## Template functions

-   `prism_path` : take a node as parameter and generates a valid path.

## Routes definition

-   `prism_format` : Generates an url like this: `/:nid.:format`
-   `prism`: Generates an url like this: `/:nid`
-   `prism_path_format`: Generates an url like this: `/:path.:format`
-   `prism_path`: Generates an url like this: `/:path.:format`
