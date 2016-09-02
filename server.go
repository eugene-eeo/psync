package main

import (
	"github.com/gin-gonic/gin"
)

func Serve(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/", PsyncBlocksDir())
	r.Run(addr)
}
