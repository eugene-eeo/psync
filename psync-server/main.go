package main

import (
	"github.com/oxtoacart/bpool"
	"github.com/gin-gonic/gin"
	"os"
	"os/user"
	"path/filepath"
)

func PsyncBlocksDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".psync/blocks")
}

func InitHome() {
	os.MkdirAll(PsyncBlocksDir(), 0755)
}

func main() {
	InitHome()
	args := os.Args[1:]
	if len(args) != 1 {
		print("usage: psync-server <addr>\n")
		os.Exit(1)
	}
	addr := args[0]
	root := PsyncBlocksDir()
	bufpool := bpool.NewBytePool(20, 4096)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:checksum", func(c *gin.Context) {
		checksum := c.Param("checksum")
		blobpath := filepath.Join(root, checksum)
		f, err := os.Open(blobpath)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}
		defer f.Close()
		buff := bufpool.Get()
		defer bufpool.Put(buff)
		f.Read(buff)
		c.Writer.Write(buff)
	})
	r.Run(addr)
}
