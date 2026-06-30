package query

import (
	"net/http"
	"strconv"
	"yoink/ranking"
	utils "yoink/utils/error"

	"github.com/labstack/echo/v5"
)

func Ranking(c *echo.Context) error{
	q := c.QueryParam("q")
	k := c.QueryParam("k")
	kInt, err := strconv.Atoi(k)
	if err != nil{
		return utils.NewApiError(c, http.StatusBadRequest, err.Error())
	}

	pages, err := ranking.RankPages(q, kInt)
	if err != nil{
		return utils.NewApiError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pages)
}