package crawl

import (
	"fmt"
	"net/http"
	"yoink/models"
	"yoink/utils/database"
	utils "yoink/utils/error"

	"github.com/labstack/echo/v5"
)

func NoOfPagesCrawled(c *echo.Context) error{
	db := database.DB

	var count int64
	if err := db.Model(&models.Page{}).Count(&count).Error; err != nil{
		utils.NewApiError(c, http.StatusInternalServerError, fmt.Sprintln("An internal server error has occurred"))
		return err
	}
	c.JSON(http.StatusOK, map[string]int64{
		"count": count,
	})
	return nil
}