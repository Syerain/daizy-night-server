package config

import (
	"log/slog"
	"net"

	"github.com/spf13/viper"
)

var globalConfig *Config

func GetConfig() *Config {
	return globalConfig
}

// currently stated the returned error but didnt use.
func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Failed to read Config file. Using default config values.")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	globalConfig = &cfg
	return nil
}

func setDefaults() {
	viper.SetDefault("Http.ListenPort", 4703)
	viper.SetDefault("Http.ListenAddress", "localhost")
	viper.SetDefault("Http.VersionMask", 0)
	viper.SetDefault("Main.IsDebugMode", true)
	viper.SetDefault("Database.IsDebugMode", false)
	viper.SetDefault("Log.IsDebugMode", false)
}

// havent finished yet.
func validateConfig(cfg *Config) error {
	errorCounts := 0

	// used
	switch 1 {
	case 1:
		if cfg.Http.ListenPort < 1 || cfg.Http.ListenPort > 65535 {
			errorCounts++
			slog.Error("Errors in Config:Invalid listen port. Must be between 1 and 65535.")
		}
		fallthrough
	case 2:
		parsed := net.ParseIP(cfg.Http.ListenAddress)
		if parsed == nil {
			errorCounts++
			slog.Error("Errors in Config:Invalid listen address.")
		}
		fallthrough
	case 3:

	}
	return nil
}
