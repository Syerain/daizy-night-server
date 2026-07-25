package config

import "time"

const VersionMask int = 0

// Config struct holds the configuration for the server
/* Notice that server configuration differs from client's despite they share the same name. */
type Config struct {
	Main struct {
		IsDebugMode bool `mapstructure:"isDebugMode"`
	} `mapstructure:"main"`

	Http struct {
		Port    int    `mapstructure:"port"`
		Address string `mapstructure:"address"`

		// use format like 1m30s and viper can read it
		JwtAccessTokenExpireTime  time.Duration `mapstructure:"jwtAccessTokenExpireTime"`
		JwtRefreshTokenExpireTime time.Duration `mapstructure:"jwtRefreshTokenExpireTime"`
	} `mapstructure:"http"`

	Database struct {
		IsDebugMode bool   `mapstructure:"isDebugMode"`
		DSN         string `mapstructure:"DSN"`
	} `mapstructure:"database"`

	Log struct {
		// log level is determined by MainConfig, instead of LogConfig.
		// LogConifg controls those exact style params of logs.
		IsColored bool `mapstructure:"isColored"`
	} `mapstructure:"log"`

	Security struct {
		RegistercodeEnckey    string `mapstructure:"registercodeEnckey"`
		RegistercodeDeckey    string `mapstructure:"registercodeDeckey"`
		PasswordEnckey        string `mapstructure:"passwordEnckey"`
		PasswordDeckey        string `mapstructure:"passwordDeckey"`
		JwtAccessTokenEnckey  string `mapstructure:"jwtAccessTokenEnckey"`
		JwtAccessTokenDeckey  string `mapstructure:"jwtAccessTokenDeckey"`
		JwtRefreshTokenEnckey string `mapstructure:"jwtRefreshTokenEnckey"`
		JwtRefreshTokenDeckey string `mapstructure:"jwtRefreshTokenDeckey"`
	} `mapstructure:"security"`
}
