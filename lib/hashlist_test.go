package lib_test

import (
	"github.com/eugene-eeo/psync/lib"
	"testing"
	"bytes"
	"strings"
)

var hashlist []string = []string{
	"12e489133466db18e988e4464e4f7b0993149ff214c11e34c04faf92f11e72de",
	"f48021e273eb257daaed541287b16a498d1db1e6f8d6c8b896164ae99b984b20",
	"cf20df963f4d49c98bb9c1229a259049a083dfca4fb112fba6193f8ce92293b3",
	"86180974b09ee52e5edd084476bdb0662d208f5107b6b44691b8221569bd5063",
	"6825baa392f05d5bccb1de81a7ac4763ffb612c43fa89c3339e4d9e4c07611ac",
	"",
	"29190d0155be19a3dfc7ec851742a8f200150a4ff1cb7b68c82dcbe322f6df08",
}
var hashlist_str string = strings.Join(hashlist, "\n")

func TestParseHashList(t *testing.T) {
	i := 0
	b := bytes.NewBufferString(hashlist_str)
	lib.ParseHashList(b, func(line lib.Checksum) {
		if string(line) != hashlist[i] {
			t.Error("expected", line, "to match", hashlist[i])
		}
		i++
	})
	if i != 5 {
		t.Error("expected only 5 hashes, got", i)
	}
}

func TestNewHashList(t *testing.T) {
	b := bytes.NewBufferString(hashlist_str)
	hl := lib.NewHashList(b)
	if len(*hl) != 5 {
		t.Error("expected hashlist to have length 5")
	}
}
