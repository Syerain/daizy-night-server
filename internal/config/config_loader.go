package config

import "github.com/atomreforge/confx"

// MustLoadConfig panics when encountering errors.
func MustLoadConfig() *Config {
	cfg := confx.MustLoad[Config]("daizynight",
		//conf.WithLocale("zh"),

		// Explicitly specified
		confx.WithFileName("config"),
		confx.WithFileType("yaml"),
	)

	return cfg
}
