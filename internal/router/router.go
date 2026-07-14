package router

import (
	"daizynight/internal/handler"

	"github.com/labstack/echo/v5"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/api/register", handler.HandleRegister)

	return e
}
