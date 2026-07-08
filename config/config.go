package config

// Config struct holds the configuration for the server
// Notice that server configuration differs from client's despite they share the same name.
type Config struct {
	Http     HttpConfig     `mapstructure:"httpConfig"`
	Database DatabaseConfig `mapstructure:"databaseConfig"`
	Log      LogConfig      `mapstructure:"logConfig"`
}

type HttpConfig struct {
	ListenPort    int    `mapstructure:"listenPort"`
	ListenAddress string `mapstructure:"listenAddress"`
	// versionMask is used to tell the server which version of request is appropriate to use.
	// a request with a version number smaller than the versionMask will be rejected by the server.
	VersionMask int `mapstructure:"versionMask"`
}

type DatabaseConfig struct {
	IsDebugMode bool `mapstructure:"isDebugMode"`
}

type LogConfig struct {
	IsDebugMode bool `mapstructure:"isDebugMode"`
}
