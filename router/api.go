package router

import (
	"yoink/controller"
	"yoink/controller/query"

	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo){
	api := e.Group("/api")
	api.GET("/",  controller.API)

	api.GET("/query", query.Ranking)
}