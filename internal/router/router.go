package router

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger = slog.Default()
	return e
}
