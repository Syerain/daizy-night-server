package main

import (
	"daizynight/config"
	"fmt"
	"log/slog"
)

// at first I want to use custom loggers however there's too much work.
// maybe I will use them in the future but not for now.
/* func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)

	log = slog.New(handler)
} */

func main() {

	slog.Info("server starting...")
	slog.Info("Reading config...")

	config.InitConfig()
	cfg := config.GetConfig()
	testnum := cfg.Http.ListenPort
	fmt.Println("config listen port: ", testnum)

	slog.Info("config listen address: ", cfg.Http.ListenAddress)
}
