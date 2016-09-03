package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/eugene-eeo/psync/lib"
)

func Serve(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/", lib.PsyncBlocksDir())
	r.Run(addr)
}
