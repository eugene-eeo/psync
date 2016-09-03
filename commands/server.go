package commands

import (
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/eugene-eeo/psync/lib"
	"github.com/oxtoacart/bpool"
)

func Serve(addr string) {
	root := lib.PsyncBlocksDir()
	pool := bpool.NewBytePool(20, lib.BLOCK_SIZE)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:checksum", func(c *gin.Context) {
		checksum := c.Param("checksum")
		f, err := os.Open(filepath.Join(root, checksum))
		if err != nil {
			c.AbortWithStatus(404)
		}
		defer f.Close()
		buff := pool.Get()
		f.Read(buff)
		c.Writer.Write(buff)
		pool.Put(buff)
	})
	r.Run(addr)
}
