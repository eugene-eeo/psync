# psync

p2p block distribution service thing. you can choose to run a
server that serves some blocks (4096-byte chunks):

    $ psync-server export image.png | tee image.png.hashlist
    ba7816bf...
    cb8379ac...
    $ psync-server up localhost:8000

and then others will get data from your servers if they trust you:

    $ pysnc get localhost:8000 image.png.hashlist > image.png.2
    $ diff image.png image.png.2

there will be no mechanism for hashlist distribution, because you and
I may differ on what `image.png` means. If I manage to convince myself
to a solution, maybe a checksum of hashlists then I will implement it.

## todo

 - use sqlite for storage
 - come up with method for peer sharing
 - split into separate tools
