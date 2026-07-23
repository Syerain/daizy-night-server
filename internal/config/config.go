package config

// Config struct holds the configuration for the server
/* Notice that server configuration differs from client's despite they share the same name. */
type Config struct {
	Main     MainConfig     `mapstructure:"Main"`
	Http     HttpConfig     `mapstructure:"Http"`
	Database DatabaseConfig `mapstructure:"Database"`
	Log      LogConfig      `mapstructure:"Log"`
	Security SecurityConfig `mapstructure:"Security"`
}

type MainConfig struct {
	IsDebugMode bool `mapstructure:"IsDebugMode"`
}
type HttpConfig struct {
	ListenPort    int    `mapstructure:"ListenPort"`
	ListenAddress string `mapstructure:"ListenAddress"`
	/* VersionMask is used to tell the server which versions of client request are appropriate to use.
	a request with a version number smaller than the VersionMask will be rejected by the server. */
	VersionMask int `mapstructure:"VersionMask"`
}

type DatabaseConfig struct {
	IsDebugMode bool   `mapstructure:"IsDebugMode"`
	DSN         string `mapstructure:"DSN"`
}

type LogConfig struct {
	// log level is determined by MainConfig, instead of LogConfig.
	// LogConifg controls those exact style params of logs.
	IsColored bool `mapstructure:"IsColored"`
}

type SecurityConfig struct {
	RegistercodeEnckey string `mapstructure:"RegistercodeEnckey"`
	RegistercodeDeckey string `mapstructure:"RegistercodeDeckey"`
	AccessTokenEnckey  string `mapstructure:"AccessTokenEnckey"`
	AccessTokenDeckey  string `mapstructure:"AccessTokenDeckey"`
}
