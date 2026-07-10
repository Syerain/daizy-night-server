package config

import (
	"errors"
	"log/slog"
	"net"

	"github.com/spf13/viper"
)

// get it via GetConfig()
var globalConfig *Config

func GetConfig() *Config {
	return globalConfig
}

// currently stated the returned error but didnt use.
func InitConfig() error {
	// viper params
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// viper setting default values.
	setDefaults()

	// trying to read in config file and validate then.
	// any errors will be throwed upward.
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Failed to read Config file. Using default config values.")
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	if err := validateConfig(&cfg); err != nil {
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
	viper.SetDefault("Log.IsColored", true)
}

// havent finished yet.
func validateConfig(cfg *Config) error {
	var errs []error

	if cfg.Http.ListenPort < 1 || cfg.Http.ListenPort > 65535 {
		errs = append(errs, &ConfigValidationError{
			Field:   "Http.ListenPort",
			Message: "must be between 1 and 65535",
		})
	}
	if parsedIP := net.ParseIP(cfg.Http.ListenAddress); parsedIP == nil {
		errs = append(errs, &ConfigValidationError{
			Field:   "Http.ListenAddress",
			Message: "must be a valid IP address",
		})
	}
	if cfg.Http.VersionMask < 0 {
		errs = append(errs, &ConfigValidationError{
			Field:   "Http.VersionMask",
			Message: "must be a non-negative integer",
		})
	}
	return errors.Join(errs...)
}
