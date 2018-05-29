package main

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/middleware/logger"
)

func main()  {
    app := iris.New()

    customLogger := logger.New(logger.Config{
        Status: true,
        IP: true,
        Method: true,
        Path: true,
        MessageContextKeys: []string{"logger_message"},
        MessageHeaderKeys: []string{"User-Agent"},
    })
    app.Use(customLogger)
    app.Logger().SetLevel("debug")
    app.Get("/", func(ctx iris.Context) {
        ctx.HTML("<b>Hello!</b>")
    })

    app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
}
