package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/labstack/echo/v5"
)

func CallChainPush(ctx *echo.Context, step string) func() {
	stack, _ := echo.ContextGet[[]string](ctx, string(consts.JsonExprCallChain))
	stack = append(stack, step)
	ctx.Set(string(consts.JsonExprCallChain), stack)

	return func() {
		cur, _ := echo.ContextGet[[]string](ctx, string(consts.JsonExprCallChain))
		if len(cur) > 0 {
			ctx.Set(string(consts.JsonExprCallChain), cur[:len(cur)-1])
		}
	}
}

// AppendCallChain faithfully records each step a request passes through
// (monotonically accumulating, never popped). Differs from CallChainPush:
// it returns nothing and is meant for the error-exit to show the full flow
// chain as it actually happened.
func AppendCallChain(ctx *echo.Context, step string) {
	stack, _ := echo.ContextGet[[]string](ctx, string(consts.JsonExprCallChain))
	stack = append(stack, step)
	ctx.Set(string(consts.JsonExprCallChain), stack)
}


// summon random traceid
func NewTraceID() string {
	i := make([]byte, 8)
	rand.Read(i)
	return hex.EncodeToString(i)
}

func GetCallChainCurrentStep(ctx *echo.Context) string {
	stack, _ := ctx.Get(string(consts.JsonExprCallChain)).([]string)
	if len(stack) == 0 {
		return ""
	}
	return stack[len(stack)-1]
}

func GetCallChain(ctx *echo.Context) []string {
	stack, _ := ctx.Get(string(consts.JsonExprCallChain)).([]string)
	return stack
}

// used before a slog.Info() to output more custom info
func Layer(ctx *echo.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	tid, _ := ctx.Get(string(consts.ExprTraceid)).(string)
	step := GetCallChainCurrentStep(ctx)

	l := slog.Default().With(slog.String(string(consts.ExprTraceid), tid))
	if step != "" {
		l = l.With(slog.String("step", step))
	}
	return l
}
