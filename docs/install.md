Installation
============

Please note, this is a work in progress. There is no real installation for a custom project. For now, you have to work
from wihin the current repository.

Requirements
------------

- Backend: You must have GO 1.4+ installed, and a running instance of PostgreSQL running.
- Frontend: You must have ``nodejs`` and ``npm`` installed

Installation steps
------------------

1. Retrieve the code source: ``go get github.com/rande/gonode/core``
2. Configure the ``server.toml`` configuration file
3. Start the webserver: ``make run``
4. Create a valid schema: ``curl -XPOST http://localhost:2508/setup/install`` 
5. Load some fixtures: ``curl -XPOST http://localhost:2508/setup/data/load`` 
 
DX
--

* ``git update-index --assume-unchanged app/assets/bindata.go``