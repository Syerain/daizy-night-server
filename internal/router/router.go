package router

import "github.com/labstack/echo/v5"

func New() *echo.Echo {
	e := echo.New()
	return e
}
