package middleware

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/labstack/echo/v5"
)

func RespondCustom(ctx *echo.Context, errapp abstract.InterfaceAppError) error {
	return ctx.JSON(errapp.StatusCode(),
		map[string]string{"message": string(errapp.Error())})
}

func Respond(ctx *echo.Context, status int, msg string) error {
	return ctx.JSON(status, map[string]string{"message": msg})
}

func RespondObj(ctx *echo.Context, status int, obj any) error {
	return ctx.JSON(status, obj)
}
