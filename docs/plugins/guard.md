Guard
=====

Introduction
------------

Guard plugin handle request authentification. The plugins comes with a dedicated middleware and request authenticators.
For now there is only 2 authenticators implemented:
 - ``JwtLoginGuardAuthenticator``: create a valid Json Web Token from the posted ``username`` and ``password``.
 - ``JwtTokenGuardAuthenticator``: validate a Json Web Token. 

Configuration
-------------


    ```toml
    [guard]
    key = "ZeSecretKey0oo"
    
        [guard.jwt]
            [guard.jwt.login]
            path = "/login"
    
            [guard.jwt.token]
            path = "^\\/nodes\\/(.*)$"


- ``key`` is private and it is used to sign the JWT with a symetric algorythm.
- ``guard.jwt.login.path`` is used to configure the login entry point, ie where the ``JwtLoginGuardAuthenticator`` will accept the request.
- ``guard.jwt.token.path`` is used to configure paths requiring to have authentification handled by the ``JwtTokenGuardAuthenticator`` service.


Authenticators
--------------

### JwtLoginGuardAuthenticator


The service will use the ``core.user`` node type to find the user by her/his username. The query looks like: ``type = 'core.user' AND data->>'username' = ?``

The authentification request should be a POST

```HTTP
POST /login HTTP/1.1
Content-Type: application/x-www-form-urlencoded

username=admin&password=secret
```

If the response is valid, the response will be:

```HTTP
HTTP/1.1 200 OK
Content-Type: application/json

{
    "status":  "OK",
    "message": "Request is authenticated",
    "token":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTA0Nzg1NzQsInJscyI6bnVsbCwidXNyIjoicmFuZGUifQ.E_BMRg2UWO7jVw1CGgn7WhhwbATCHjYYcausZZ7LSZA",
}

```

If the response is not valid, the response will be

```HTTP
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
    "status":  "KO",
    "message": "Unable to authenticate request"
}
```
   
### JwtTokenGuardAuthenticator

The service will use the ``core.user`` node type to find the user by her/his username. The query looks like: ``type = 'core.user' AND data->>'username' = ?``

The authentification request should be on any http method, either using the ``Authorization`` header or the ``access_token`` parameter.

```HTTP
GET /nodes HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTA0Nzg1NzQsInJscyI6bnVsbCwidXNyIjoicmFuZGUifQ.E_BMRg2UWO7jVw1CGgn7WhhwbATCHjYYcausZZ7LSZA
```

or

```HTTP
GET /nodes?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTA0Nzg1NzQsInJscyI6bnVsbCwidXNyIjoicmFuZGUifQ.E_BMRg2UWO7jVw1CGgn7WhhwbATCHjYYcausZZ7LSZA HTTP/1.1
```
