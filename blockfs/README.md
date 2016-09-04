# blockfs

A very simple layer over hashlists and blocks. Manages a simple directory of
blocks (2MB chunks of data+checksum) and named hashlists. Usage example for
storing an external image into the store, saving the checksums into a hashlist,
and piping the chunks into stdout:

```
import "os"
import "github.com/eugene-eeo/psync/blockfs"

func main() {
    f, _ := os.Open("/external/image.png")
    fs := blockfs.NewFS(".")
    hashlist, _ := fs.ExportNamed(f, "my-image")
    for _, checksum := range hashlist {
        block, err := fs.GetBlock(checksum)
        if err != nil {
            block.WriteTo(os.Stdout)
        }
    }
}
```

## todo

 - lots of tests
 - implement extra features:
    - gc for cleaning up unused blocks
