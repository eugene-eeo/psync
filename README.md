# psync

p2p block distribution service thing. you can choose to run a server
that serves some files (normalised to hash-lists):

    $ psync export image.png | tee image.png.hashlist
    ba7816bf...
    cb8379ac...
    $ psync server localhost:8000

and then others will get data from your servers:

    $ pysnc get localhost:8000 image.png.hashlist > image.png.2
    $ diff image.png image.png.2
