package config

// Config struct holds the configuration for the server
// Notice that server configuration differs from client's despite they share the same name.
type Config struct {
	Main     MainConfig     `mapstructure:"Main"`
	Http     HttpConfig     `mapstructure:"Http"`
	Database DatabaseConfig `mapstructure:"Database"`
	Log      LogConfig      `mapstructure:"Log"`
}

type MainConfig struct {
	IsDebugMode bool `mapstructure:"IsDebugMode"`
}
type HttpConfig struct {
	ListenPort    int    `mapstructure:"ListenPort"`
	ListenAddress string `mapstructure:"ListenAddress"`
	// VersionMask is used to tell the server which version of request is appropriate to use.
	// a request with a version number smaller than the versionMask will be rejected by the server.
	VersionMask int `mapstructure:"VersionMask"`
}

type DatabaseConfig struct {
	IsDebugMode bool `mapstructure:"IsDebugMode"`
}

type LogConfig struct {
	isColored bool `mapstructure:"IsColored"`
}
