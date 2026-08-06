package handler

import (
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/utils"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleAdminSudo(ctx *echo.Context) error {
	// record flow chain (monotonically accumulating)
	utils.AppendCallChain(ctx, string(consts.ModExprHandlerAdmin))
	utils.AppendCallChain(ctx, string(consts.ModExprServiceAdmin))

	h.ServiceAdmin.Sudo()
	return nil
}
