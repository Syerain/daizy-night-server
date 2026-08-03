package handler

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/labstack/echo/v5"
)

/*
// THIS FUNC IS ABANDONED !
func ErrorRespondHttp(ctx *echo.Context, err error) error {
	return nil
	var errapp abstract.InterfaceAppError
	if errors.As(err, &errapp) {
		return ctx.JSON(errapp.HttpAbort(), map[string]string{"message": errapp.Error()})
	}
	unknown := errs.ErrUnknown{
		Type: errs.Unknown,
		Http: 520,
	}
	slog.Error(unknown.Type.Say())
	return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "internal error"})
} */

func RespondCustom(ctx *echo.Context, errapp abstract.InterfaceAppError) error {
	return ctx.JSON(errapp.HttpAbort(),
		map[string]string{"message": string(errapp.Error())})
}

func Respond(ctx *echo.Context, status int, msg string) error {
	return ctx.JSON(status, map[string]string{"message": msg})
}

func RespondObj(ctx *echo.Context, status int, obj any) error {
	return ctx.JSON(status, obj)
}
