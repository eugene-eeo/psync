package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/eugene-eeo/psync/lib"
)

func Serve(addr string) {
	root := lib.BlocksDir()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/", root)
	CheckError(r.Run(addr))
}
