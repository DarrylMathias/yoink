package main

import (
	"yoink/app"
	"yoink/ranking"
	"yoink/router"
	"yoink/utils/env"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)


func main() {
    app.App()
	ranking.Init()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	router.RegisterRoutes(e)
	e.Logger.Error("failed to start server", "error", e.Start(env.ConfigValue.ApplicationPort)) 

}
