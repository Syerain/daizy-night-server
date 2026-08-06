package handler

import "github.com/labstack/echo/v5"

func (h *HandlerComplex) HandleAdminSudo(ctx *echo.Context) error {
	h.ServiceAdmin.Sudo()
	return nil
}
