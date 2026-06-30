package main

import (
	"fmt"
	"net/http"
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
	fmt.Println("initialized ranking server")

	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{env.ConfigValue.FrontEndURL},
		AllowMethods: []string{http.MethodGet},
	}))

	router.RegisterRoutes(e)
	e.Logger.Error("failed to start server", "error", e.Start(env.ConfigValue.ApplicationPort)) 

}
