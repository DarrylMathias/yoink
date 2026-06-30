package query

import (
	"fmt"
	"net/http"
	"yoink/ranking"
	utils "yoink/utils/error"

	"github.com/labstack/echo/v5"
)

func Ranking(c *echo.Context) error{
	fmt.Println("hit ranking")
	q := c.QueryParam("q")

	pages, err := ranking.RankPages(q)
	if err != nil{
		return utils.NewApiError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pages)
}