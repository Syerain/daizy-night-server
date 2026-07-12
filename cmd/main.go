package main

import (
	"daizynight/internal/config"
	"daizynight/internal/router"
	"daizynight/internal/utils"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// reading config, temply using slog vanilla logger.
	slog.Info("Server starting...")
	slog.Info("Reading config...")

	// initializa colored logger
	err := config.InitConfig()
	if err != nil {
		slog.Error("FATAL: Failed to initialize config", "error", err)
		return
	}
	cfg := config.GetConfig()

	// initializing our colored logger
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	// printing basic config info
	if cfg.Main.IsDebugMode {
		slog.Info("Debug mode is enabled.")
	}

	slog.Info("Config listen port: " + fmt.Sprintf("%d", cfg.Http.ListenPort))
	slog.Info("Config listen address: " + cfg.Http.ListenAddress)

	// loading Echo engine
	slog.Info("Loading HTTP server ...")

	e := router.New()
	addrport := fmt.Sprintf("%s:%d", cfg.Http.ListenAddress, cfg.Http.ListenPort)

	slog.Info("Listening on " + addrport)

	if err := e.Start(addrport); err != nil {
		slog.Error("FATAL: Failed to start HTTP server.")
		os.Exit(1)
	}
}
