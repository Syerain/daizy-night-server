package config

import "time"

const MinimalAPIVersionRequired int = 0 // compile-time config

// Config struct holds the configuration for the server
/* Notice that server configuration differs from client's despite they share the same name. */
type Config struct {
	Main struct {
		IsDebugMode bool `mapstructure:"isDebugMode" default:"true"`
	} `mapstructure:"main"`

	Http struct {
		Port    uint   `mapstructure:"port" validate:"port" default:"4703"`
		Address string `mapstructure:"address" validate:"ip_addr" default:"127.0.0.1"`
	} `mapstructure:"http"`

	Database struct {
		IsDebugMode bool   `mapstructure:"isDebugMode" default:"false"`
		DSN         string `mapstructure:"DSN" validate:"required"`
	} `mapstructure:"database"`

	Log struct {
		// log level is determined by MainConfig, instead of LogConfig.
		// LogConifg controls those exact style params of logs.
		IsColored bool `mapstructure:"isColored" default:"false"`
	} `mapstructure:"log"`

	Security struct {
		RegistercodeEnckey    string `mapstructure:"registercodeEnckey" validate:"required"`
		RegistercodeDeckey    string `mapstructure:"registercodeDeckey" validate:"required"`
		PasswordEnckey        string `mapstructure:"passwordEnckey" validate:"required"`
		PasswordDeckey        string `mapstructure:"passwordDeckey" validate:"required"`
		JwtAccessTokenEnckey  string `mapstructure:"jwtAccessTokenEnckey" validate:"required"`
		JwtAccessTokenDeckey  string `mapstructure:"jwtAccessTokenDeckey" validate:"required"`
		JwtRefreshTokenEnckey string `mapstructure:"jwtRefreshTokenEnckey" validate:"required"`
		JwtRefreshTokenDeckey string `mapstructure:"jwtRefreshTokenDeckey" validate:"required"`

		// use format like 1m30s and viper can read it
		JwtAccessTokenExpireTime  time.Duration `mapstructure:"jwtAccessTokenExpireTime" validate:"required"`
		JwtRefreshTokenExpireTime time.Duration `mapstructure:"jwtRefreshTokenExpireTime" validate:"required"`
	} `mapstructure:"security"`
}
