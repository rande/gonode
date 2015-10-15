Vault
=====

Introduction
------------

The vault is a place used to store binary files, a vault is a key/value store specialized on storing
binary content.

A vault has a few methods: Has, Get, GetReader, Put and Delete. There is no search or listing options,
these feature need to be done by the upper layer of the application. 

Concept
-------

A vault stores files inside a data container, when a file is stored up to 2 extra files are created:

 - ``Vaultfile``: the Vaultfile contains the key used to encrypt the file and the md5 file of the 
original file
 - ``Metafile``: the Metafile contains any metadata linked to the file. Depends on the vault file this 
 file can be created or not (fs vs s3)
 
The encrypted key from the vaultfile is used to encrypt the meta file and the binary file.

A vault also have a main encryption key use to encrypt the vaultfile. If no key is provided, then no
encryption will be used for the vaultfile. (ie, an attacker will be able to descrypt the binary file 
and the meta file.

The vault configuration (vault type, location and main encryption key) must be handled by the upper layers.

Name
----

The ``name`` is used to generate an unique value used to store the file. The ``name`` can be any string
value. So with the ``VaultFs`` the name value will be transformed into a sha1 hash used to generated
the final internal path:

 - name: "this-is-a-test"
 - sha1: 7b87fd8ec71a47da643cd1f06c3e6b7ef42d8554
 - binary path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin
 - metafile path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin.meta
 - vaultfile path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin.vault
 
So the name must be unique inside the store. 

Key
---

A vaultfile ``key`` is generated using the ``generateKey`` function which return a random 32 bytes array.

Encrypter/Decrypter
-------------------

A vault also have a set of encrypter/descrypted functions used to encrypt/descript the file on the fly.

There are 2 couples:

  - ``Noop`` : no operation, ie no encryption applied. This can be usefull for debugging or for 
  storing non critical information (ie, web site assets)
  - ``AES_OFB``: apply AES encryption with OFB Mode. This mode is good for only offer confidentiality, 
  there is no authenticity and integrity encryption. With this mode it is still possible to flip a 
  bit and change the decrypted contents. However, it is possible to encrypt a large stream.
  
  
Warning
-------

 - For now there is no hmac verification on contents streamed.
 - There is a need to implement AES with CBC block mode inside a io.Reader compatible interface by 
 splitting encryption/decryption into fixed chunk size block.
 - Possibly need to adapt code from: https://github.com/tadzik/simpleaes/blob/master/simpleaes.go to
 work with current interface signature.

Vault
-----

For now, there is only one vault implemented:
 
 - ``VaultFs``: use the current filestem to store file

Planned:

 - ``VaultS3``: store file into a S3 bucket
 - ``VaultDuplicate``: proxy to store file into multiple vault (cheap replication)