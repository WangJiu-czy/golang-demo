package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

const (
	host = "127.0.0.1:8080"
)

func main() {
	app := iris.New()
	rv := router.NewRoutePathReverser(app, router.WithHost(host), router.WithScheme("http"))
	templates := iris.HTML("templates", ".html")
	templates.AddFunc("url", rv.URL)
	app.RegisterView(templates)

	subdomain := app.WildcardSubdomain()

	mypathRoute := subdomain.Get("/mypath", emptyHandler)
	mypathRoute.Name = "my-page1"

	mypath2Route := subdomain.Get("/mypath2/{paramfirst}/{paramsecond}", emptyHandler)
	mypath2Route.Name = "my-page2"

	mypath3Route := subdomain.Get("/mypath3/{paramfirst}/statichere/{paramsecond}", emptyHandler)
	mypath3Route.Name = "my-page3"

	mypath4Route := subdomain.Get("/mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}", emptyHandler)
	mypath4Route.Name = "my-page4"

	mypath5Route := subdomain.Handle("GET", "/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", emptyHandler)
	mypath5Route.Name = "my-page5"

	mypath6Route := subdomain.Get("/mypath6/{paramfirst}/{paramsecond}/staticParam/{paramThirdAfterStatic}", emptyHandler)
	mypath6Route.Name = "my-page6"

	app.Get("/", func(ctx iris.Context) {
		// for username5./mypath6...
		paramsAsArray := []string{"username5", "theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.html"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})
	app.Get("/mypath7/{paramfirst}/{paramsecond}/static/{paramthird}", emptyHandler).Name = "my-page7"
	app.Listen(host)
}
func emptyHandler(ctx iris.Context) {
	ctx.Writef("Hello from subdomain: %s , you're in path:  %s", ctx.Subdomain(), ctx.Path())
}
