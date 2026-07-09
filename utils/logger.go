package utils

import (
	"log/slog"
	"os"
	"time"

	"github.com/pwntr/tinter"
)

// scheduled logs with module name and debug level but havnt finished yet.
// to regulate the identification,
// moduleName stated here is string but in fact you shall pass an string-based enum value.
func InitModuleLogger(isDebugMode bool, moduleName string) {
	var handler slog.Handler

	level := slog.LevelInfo
	if isDebugMode {
		level = slog.LevelDebug
	}

	handler = tinter.NewHandler(os.Stdout, &tinter.Options{
		Level:      level,
		TimeFormat: time.RFC3339,
	})

	logger := slog.New(handler).With(
		slog.String("module", moduleName),
	)

	slog.SetDefault(logger)
}

/*
type moduleHandler struct {
	handler slog.Handler
	name    string
} */
