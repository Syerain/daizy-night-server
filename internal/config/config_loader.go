package config

import (
	"github.com/oy3o/conf"
)

// MustLoadConfig panics when encountering errors.
func MustLoadConfig() *Config {
	cfg := conf.MustLoad[Config]("daizy",
		conf.WithLocale("zh"),

		// Explicitly specified
		conf.WithFileName("config"),
		conf.WithFileType("yaml"),
	)

	return cfg
}
