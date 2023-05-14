package main

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

func main() {

	app := iris.New()
	app.Get("/getJWT", func(ctx iris.Context) {
		// 往jwt中写入了一对值,用户登录成功后,将信息写入并生成token
		token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
		})

		// 签名生成jwt字符串
		tokenString, _ := token.SignedString([]byte("czy"))

		// 返回
		ctx.JSON(tokenString)
	})
	app.Get("/showHello", j.Serve, func(c iris.Context) {
		jwtInfo := c.Values().Get("jwt").(*jwt.Token)

		c.JSON(jwtInfo)
	})
	app.Run(iris.Addr(":8080"))
}
