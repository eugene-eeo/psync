# psync

p2p block distribution service thing. you can choose to run a server
that serves some files (normalised to hash-lists):

    $ psync export image.png | tee image.png.hashlist
    ba7816bf...
    cb8379ac...
    $ psync server localhost:8000 &
    $ pysnc get localhost:8000 image.png.hashlist > image.png.2
    Fetching blocks:
     - ba7816bf...
     - cb8379ac...
    Done
    $ diff image.png image.png.2
