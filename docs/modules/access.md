Access
======

Introduction
------------

The access module is part of the gonode package, it is responsible of handling security at the http api and node level. 

Configuration
-------------

_No configuration options available_


Security Workflow
-----------------

There are 3 securities levels, there are all using the security token to check access permissions:

- The firewall: check token's roles against an http request using regular expression. 
- The action: ckeck token's roles against a predefined set in the api or prism
- The node: check token's roles against the ``access`` field. 

Action Access
-------------

The different actions are protected by roles as defined in [restful_api.md](restful_api.md).

Node Access
-----------

Any actions perfomed on a node is controlled by the ``access`` field. If one value of the ``access`` is present in the security token, then the action is authorized.
 
In case of a unit action (create, delete, update) a 403 error will be raised, however for listing, the result will be filtered.

Please note: the ``node:api:master`` role will allow any actions to be performed.