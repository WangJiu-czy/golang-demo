package main

import (
	"github.com/kataras/iris/v12"
)

type mypage struct {
	Title   string
	Message string
}

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("templates", ".html").Layout("layout.html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.CompressWriter(true)
		ctx.ViewData("", mypage{Title: "My Page title", Message: "Hello world"})
		if err := ctx.View("mypage.html"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}

	})
	app.Listen(":8080")
}
