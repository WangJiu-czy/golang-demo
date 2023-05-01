package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	tmpl := iris.HTML("templates", ".html")
	tmpl.Layout("layouts/layout.html")
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})
	app.RegisterView(tmpl)
	app.Get("/", func(c iris.Context) {
		if err := c.View("page1.html"); err != nil {

			c.StatusCode(iris.StatusInternalServerError)
			c.Writef(err.Error())
		}
	})

	app.Get("/nolayout", func(c iris.Context) {
		c.ViewLayout(iris.NoLayout)
		if err := c.View("page1.html"); err != nil {
			c.StatusCode(iris.StatusInternalServerError)
			c.Writef(err.Error())
		}
	})

	my := app.Party("/my").Layout("layouts/mylayout.html")
	{ // both of these will use the layouts/mylayout.html as their layout.
		my.Get("/", func(ctx iris.Context) {
			if err := ctx.View("page1.html"); err != nil {
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
		})
		my.Get("/other", func(ctx iris.Context) {
			if err := ctx.View("page1.html"); err != nil {
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
		})
	}
	app.Listen(":8080")
}
