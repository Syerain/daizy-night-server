package handler

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
