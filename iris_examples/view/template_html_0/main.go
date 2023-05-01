package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings" + s + "!"
	})
	app.RegisterView(tmpl)
	app.Get("/", hi)
	app.Listen(":8080", iris.WithCharset("utf-8"))

}

func hi(ctx iris.Context) {

	ctx.ViewData("Title", "Hi Page")
	ctx.ViewData("Name", "WangJiu")
	if err := ctx.View("hi.html"); err != nil {
		ctx.HTML("<h3>%s<h3>", err.Error())
		return
	}

}
