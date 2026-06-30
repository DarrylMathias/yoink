package query

import (
	"fmt"
	"net/http"
	"strconv"
	"yoink/ranking"
	utils "yoink/utils/error"

	"github.com/labstack/echo/v5"
)

func Ranking(c *echo.Context) error{
	q := c.QueryParam("q")
	if q == ""{
		return utils.NewApiError(c, http.StatusBadRequest, fmt.Sprintln("q is required"))
	}
	if len(q) > 100{
		return utils.NewApiError(c, http.StatusBadRequest, fmt.Sprintln("q is too long"))
	}

	k := c.QueryParam("k")
	kInt, err := strconv.Atoi(k)
	if err != nil{
		c.Logger().Error(err.Error())
		return utils.NewApiError(c, http.StatusBadRequest, fmt.Sprintln("An internal server error has occurred"))
	}
	if kInt >= 100{
		kInt = 100
	}

	pages, err := ranking.RankPages(q, kInt)
	if err != nil{
		c.Logger().Error(err.Error())
		return utils.NewApiError(c, http.StatusInternalServerError, fmt.Sprintln("An internal server error has occurred"))
	}

	return c.JSON(http.StatusOK, pages)
}