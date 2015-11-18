Contributing
============

This project has a complete test suite to avoid any issue while implementing new feature and to ensure there is no
regression. Before running the test suite, please make sure all components are correctly setup.

Project Configuration
---------------------

### Install Go

 * OSX: brew install go

## Install Node 

Node is used for frontend developement.

 * OSX: brew install node
   
### Configuring PostgreSQL

 * create a database with a specific user. All tests will be run using tables prefixed by ``test_``. Each functionnal
 test will recreate the schema.

 * Adjust he connection settings available in the ``config_test.toml`` file.
 
### Configuring AWS S3

 * Create a dedicated IAM User with a limited set of privilege. Please note this informatino are CONFIDENTIALS and should
stay on your local computer. NEVER PUSH THEM INTO GITHUB, NEVER.

 * Create a file ``~/.aws/credentials`` with the following lines
 
        [gonode-test]
        aws_access_key_id = YOUR_ACCESS_KEY_ID
        aws_secret_access_key = YOUR_SECRET_ACCESS_KEY

 * Create a bucket named ``gonode-test-GITHUB-NAME`` or any relevant name

### Environnement Variables
    
 - ``GONODE_TEST_OFFLINE``: set this variable to disable S3 tests
 - ``GONODE_TEST_AWS_VAULT_S3_BUCKET``: define the bucket name
 - ``GONODE_TEST_AWS_VAULT_ROOT``: define the root folder on the S3 bucker
 - ``GONODE_TEST_AWS_CREDENTIALS_FILE``: define the location of the credentials file
 - ``GONODE_TEST_AWS_PROFILE``: define the profile to use from the credentials file
 - ``TRAVIS``: if available the test will run extra check with large binary chunks on travisci.org

Running tests
-------------

    make test
    GONODE_TEST_OFFLINE=on make test # will not run S3 tests
    TRAVIS=on make test # will start with extra tests
    

Sending a Pull Request
----------------------

When you send a PR, just make sure that:

* You add valid test cases.
* Tests are green.
* The related documentation is up-to-date.
* Also don't forget to add a comment when you update a PR with a ping to the maintainer (``@username``), so he/she will get a notification.
