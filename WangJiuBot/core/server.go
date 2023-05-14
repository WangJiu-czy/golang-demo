package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var Server *gin.Engine

func init() {
	Server = gin.New()
}

func RunServer() {
	if WangJiu.GetBool("enable_http_server", true) == false {
		return
	}
	Server.GET("/", func(c *gin.Context) {
		c.String(200, "--WangJiu_Bot")
	})
	gin.SetMode(gin.ReleaseMode)
	logger.Info(fmt.Sprintf("开启httpserver----0.0.0.0:%s", WangJiu.Get("port", "8080")))
	Server.Run("0.0.0.0:" + WangJiu.Get("port", "8080"))

}
