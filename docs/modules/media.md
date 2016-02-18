Media
=====

Introduction
------------

Media stores information about a media, for now only image and youtube videos are supported.

Image Type
----------

The image type provides the current fields:
 
Meta:

 - ``Width``
 - ``Height``
 - ``Size``
 - ``ContentType``
 - ``Length``
 - ``Exif``: not used yet
 - ``Hash``: not used yet
 - ``SourceStatus``: SourceUrl processing status
 - ``SourceError``: SourceUrl processing result on error

Data: 

 - ``Reference``: contextual information about the media
 - ``Name``: image name
 - ``SourceUrl``: can be used to provided an url to retrieve the media
 
You can retrieve the binary through an api call by adding the ``?raw=1`` into the query string. The prism handler can
also be used, this entry point will provide crop and resize option.

 - ``http://localhost:2405/prism/bcb537ab-b349-4b3d-87e2-e43e17519af7`` => get the original image 
 - ``http://localhost:2405/prism/bcb537ab-b349-4b3d-87e2-e43e17519af7?mr=250`` => resize the media with the provided width
 - ``http://localhost:2405/prism/bcb537ab-b349-4b3d-87e2-e43e17519af7?mf=200,200`` => crop the media using a 200x200 square from
   the center of the image.
 
The image resize and crop only works with: ``jpg``, ``png`` and ``jpg`` files.
 
YouTube Type
------------

The YouTube type provided the following fields:

Meta (from youtube metadata)

 - ``Type``
 - ``Html``
 - ``Width``
 - ``Height``
 - ``Version``
 - ``Title``
 - ``ProviderName`` 
 - ``AuthorName`` 
 - ``AuthorUrl`` 
 - ``ProviderUrl`` 
 - ``ThumbnailUrl`` 
 - ``ThumbnailWidth`` 
 - ``ThumbnailHeight`` 

Data:

 - ``Vid``: The video id
 - ``Status``: processing status
 - ``Error``: processing error
 
The downloaded thumbnail will become a media type.
 
Configuration
-------------

```toml
[media]
    [media.image]
    allowed_widths = [100, 200]
    max_width = 300
```

 - ``allowed_widths``: slice of valid widths, if empty all widths are possible
 - ``max_width``: max allows width, if empty and ``allowed_width`` is empty then all widths are possible.
 
*To avoid any issues, those settings must be set and a cache layer configured*
