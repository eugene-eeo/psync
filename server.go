package main

import (
	"github.com/oxtoacart/bpool"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func Serve(addr string) {
	root := PsyncBlocksDir()
	bufpool := bpool.NewBytePool(20, 8192)
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
