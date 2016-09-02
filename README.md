# psync

p2p block distribution service thing. you can choose to run a
server that serves some blocks (8192-byte chunks):

    $ psync export image.png | tee image.png.hashlist
    ba7816bf...
    cb8379ac...
    $ psync up localhost:8000

and then others will get data from your servers if they find you:

    $ pysnc get localhost:8000 image.png.hashlist > image.png.2

there will be no mechanism for hashlist distribution, because you and
I may differ on what `image.png` means. If I manage to convince myself
to a solution, maybe a checksum of hashlists then I will implement it.

## wat

at its core psync is a tool for resolving hashlists. say you have
a file, `image.png` which you split into blocks of 8192 bytes and
hash them: (the checksums are shortened for brevity)

    ba7816bf   START OF FILE....
    cb8379ac   ......END OF FILE

the hashes on the left are used to make a hashlist (same order as
above) that you share with your friend (the hash list is much smaller
than the data). your friend wants to download the file, but not all
nodes have the chunks that make up the file or will give the correct
chunks. so he/she does the following:

1. make requests to some server for the all the chunks.

        GET http://psync.io/ba7816bf
        GET http://psync.io/cb8379ac

2. for each response, check if the checksum of the contents of the response matches the expected checksum:

        hash(http://psync.io/ba7816bf) == ba7816bf  (ok, keep chunk)
        hash(http://psync.io/cb8379ac) != cb8379ac  (something went wrong)

3. he/she can try to get the chunk from other servers or retry the request:

        hash(http://other1.io/cb8379ac) != cb8379ac  (nope)
        hash(http://other2.io/cb8379ac) != cb8379ac  (nope)
        hash(http://other3.io/cb8379ac) == cb8379ac  (ok, keep chunk)

once he/she has downloaded all the chunks he/she just needs to put
the chunks together and obtain the same data as you. the advantage
of hashlists is that retying downloads are cheaper since you don't
have to download the *entire* file again, and we can download files
from potentially untrusted sources, provided that the source of the
hashlist is reputable.


## todo

 - use sqlite for storage
 - come up with method for peer sharing
 - split into separate tools
