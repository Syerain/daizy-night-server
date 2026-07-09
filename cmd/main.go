package main

import (
	"daizynight/config"
	"daizynight/utils"
	"fmt"
	"log/slog"
)

func main() {

	slog.Info("Server starting...")
	slog.Info("Reading config...")

	config.InitConfig()
	cfg := config.GetConfig()

	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	if cfg.Main.IsDebugMode {
		slog.Info("Debug mode is enabled.")
	}
	slog.Info("Config listen port: " + fmt.Sprintf("%d", cfg.Http.ListenPort))
	slog.Info("Config listen address: " + cfg.Http.ListenAddress)

}
