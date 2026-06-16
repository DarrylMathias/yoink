package utils

import "github.com/labstack/echo/v5"

func NewApiError(c *echo.Context, statusCode int, errMessage string) error{
	return c.JSON(statusCode, 
		map[string]string{
			"error" : errMessage,
		},
	)
}