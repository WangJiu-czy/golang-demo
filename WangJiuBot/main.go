package main

import (
	"WangJiuBot/core"
	_ "WangJiuBot/im/wxmp"
	"github.com/gin-gonic/gin"
)

func main() {
	go core.RunServer()
	core.Server.Any("/test/", func(context *gin.Context) {
		context.JSONP(200, "{ok:test}")
	})
	select {}
}
