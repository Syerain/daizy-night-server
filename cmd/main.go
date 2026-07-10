package main

import (
	"daizynight/internal/config"
	"daizynight/internal/utils"
	"fmt"
	"log/slog"
)

func main() {
	// reading config, temply using slog vanilla logger.
	slog.Info("Server starting...")
	slog.Info("Reading config...")

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

}
