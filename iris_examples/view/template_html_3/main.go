package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("templates", ".html").Reload(true))
	mypathRoute := app.Get("/mypath", writePathHandler)
	mypathRoute.Name = "my-page1"
	mypath2Route := app.Get("/mypath2/{paramfirst}/{paramsecond}", writePathHandler)
	mypath2Route.Name = "my-page2"

	mypath3Route := app.Get("/mypath3/{paramfirst}/statichere/{paramsecond}", writePathHandler)
	mypath3Route.Name = "my-page3"

	mypath4Route := app.Get("/mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}", writePathHandler)
	// same as: app.Get("/mypath4/:paramfirst/statichere/:paramsecond/:otherparam/*something", writePathHandler)
	mypath4Route.Name = "my-page4"

	// same with Handle/Func
	mypath5Route := app.Handle("GET", "/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", writePathHandler)
	mypath5Route.Name = "my-page5"

	mypath6Route := app.Get("/mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}", writePathHandler)
	mypath6Route.Name = "my-page6"
	app.Get("/", func(ctx iris.Context) {

		paramsAsArray := []string{"theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.html"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})

	app.Get("/redirect/{namedRoute}", func(ctx iris.Context) {
		routeName := ctx.Params().Get("namedRoute")
		r := app.GetRoute(routeName)
		if r == nil {
			ctx.StatusCode(404)
			ctx.Writef("Route with name %s not found", routeName)
			return
		}

		println("The path of " + routeName + "is: " + r.Path)

		ctx.Redirect(r.Path)

	})
	app.Listen(":8080")

}

func writePathHandler(c iris.Context) {
	c.Writef("Hello from %s.", c.Path())
}
