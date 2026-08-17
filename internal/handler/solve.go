package handler

import "github.com/labstack/echo/v5"

func Bind[T any](ctx *echo.Context) (T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return b, err
	}
	return b, nil
}
