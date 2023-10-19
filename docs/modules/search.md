# Search

## Introduction

The `search` plugin provided a set of form to lookup node inside the datastore. The plugin is used to by the `api`
plugin to filter results.

## Configuration

    ```toml
    [search]
        max_result = 128

-   `max_result` set the limit of returned results in one query.

## Search filters

-   `page`: current page
-   `per_page`: number of result per page
-   `order_by`: field to order
-   `nid`: array of nid
-   `type`: array of type
-   `name`: name to filter
-   `slug`: slug to filter
-   `data.key`: array of value
-   `meta.key`: array of value
-   `status`: array of status
-   `weight`: array of weight
-   `revision`: revision number
-   `enabled`: boolean (f/0/false or t/1/true)
-   `deleted`: boolean (f/0/false or t/1/true)
-   `current`: boolean (f/0/false or t/1/true)
-   `updated_by`: array of nid
-   `created_by`: array of nid
-   `parent_nid`: array of nid
-   `set_nid`: array of nid
-   `source`: array of nid

## search.index node

This type can be used to configure an entry point with a pre-filtered index. As an example:

    ```yaml
    type: search.index
    name: Blog archives
    data:
        type: core.post

The data field accept all search filters.
