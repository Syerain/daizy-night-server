package router

import (
	"github.com/atomreforge/daizy-night-server/internal/handler"

	"github.com/labstack/echo/v5"
)

func New(h *handler.HandlerComplex) *echo.Echo {
	e := echo.New()

	e.POST("/api/register", h.HandleRegister)

	return e
}
