package commands

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func Serve(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/", filepath.Join(fs.Path, "blocks"))
	CheckError(r.Run(addr))
}
