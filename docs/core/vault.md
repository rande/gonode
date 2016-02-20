Vault
=====

Introduction
------------

The vault is a place used to store binary files, a vault is a key/value store specialized on storing
binary content.

A vault has a few public methods: Has, GetMeta, Get, Put and Remove. There is no search or listing options,
these feature need to be done by the upper layer of the application. 

A vault is linked to a driver which copy stream to a dedicated backend.

Concept
-------

A vault stores files inside a data container, when a file is stored up to 2 extra files are created:

 - ``Metafile``: the Metafile contains any metadata linked to the file. 
 - ``Vaultfile``: the Vaultfile contains the keys used to encrypt the Metafile (``MetaKey``) and the original binary file (``BinKey``)
 
The encryption keys from the ``Vaultfile`` is used to encrypt the meta file and the binary file.

A vault also have a main encryption key use to encrypt the ``Vaultfile``. If no key is provided, then no
encryption will be used for the vaultfile. (ie, an attacker will be able to decrypt the binary file 
and the meta file.)

The vault configuration (vault type, location and main encryption key storage) must be handled by the upper layers.

Name
----

The ``name`` is used to generate an unique value used to store the file. The ``name`` can be any string
value. So with the ``VaultFs`` the name value will be transformed into a sha256 hash used to generated
the final internal path:

 - name: "this-is-a-test"
 - sha256: 7b87fd8ec71a47da643cd1f06c3e6b7ef42d8554
 - binary path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin
 - metafile path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin.meta
 - vaultfile path: /tmp/goapp/test/vault/7b/87/fd8ec71a47da643cd1f06c3e6b7ef42d8554.bin.vault
 
So the name must be unique inside the store. Also the Root dir is ``/tmp/goapp/test/vault`` in this sample.

Key
---

A vaultfile ``key`` is generated using the ``generateKey`` function which return a random 32 bytes array.

Encrypter/Decrypter
-------------------

A vault also have a set of encrypter/descrypter functions used to manipulate the file on the fly.

There are 5 options:

  - ``no_op`` : no operation, ie no encryption applied. This can be usefull for debugging or for 
  storing non critical information (ie, web site assets)
  - ``aes_ofb``: apply AES encryption with OFB Mode.
  - ``aes_ctr``: apply AES encryption with CTR Mode. 
  - ``aes_cbc``: apply AES encryption with CBC Mode.
  - ``aes_gcm``: apply AES encryption with GCM Mode.
  
Please note: ``aes_ofb``, ``aes_ctr``, ``aes_cbc`` are good for confidentiality however there is no 
authenticity and integrity encryption. Please read [Block Cipher Mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation)
for mor information. This can be a solution if you need to encrypt a stream of bytes.

``aes_gcm`` is the best choice, however it will require to have enough memory to crypt and decrypt file. So for a 1GB 
file, you will need 2GB of free memory available on your system.

Vault
-----

For now, there is only 2 vaults implementation:
 
 - ``VaultFs``: use the current filestem to store file
 - ``VaultS3``: store file into a S3 bucket

Planned:

 - ``VaultDuplicate``: proxy to store file into multiple vault (cheap replication)