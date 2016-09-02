"""
usage:
    psync get <addr> <hashlist>
"""

from docopt import docopt
from hashlib import sha256
import sys
import grequests
import os


class Block:
    def __init__(self, checksum, contents):
        self.checksum = checksum
        self.contents = contents

    @classmethod
    def from_bytes(cls, data):
        return cls(sha256(data).hexdigest(), data)


def get_checksums(source):
    for checksum in source:
        checksum = checksum.strip()
        if not checksum:
            pass
        yield checksum


def get(args):
    os.makedirs('psync-blocks', exist_ok=True)
    reqs = []
    chunks = {}
    base_address = args['<addr>']
    checksums = list(get_checksums(open(args['<hashlist>'], 'r')))
    for checksum in checksums:
        path = os.path.join('psync-blocks', checksum)
        if os.path.exists(path):
            chunks[checksum] = path
            continue
        url = 'http://' + os.path.join(base_address, checksum)
        reqs.append(grequests.get(url, stream=True))

    for res in grequests.imap(reqs):
        checksum = res.request.url.rsplit('/', 1)[1]
        path = os.path.join('psync-blocks', checksum)
        data = res.content
        if Block.from_bytes(data).checksum != checksum:
            raise ValueError('invalid block: %s' % checksum)
        with open(path, 'wb') as fp:
            fp.write(data)
        chunks[checksum] = path

    for item in checksums:
        with open(chunks[item], 'rb') as fp:
            sys.stdout.buffer.write(fp.read())


if __name__ == '__main__':
    args = docopt(__doc__)
    if args['get']:
        get(args)
