package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/atomreforge/confx"
)

// MustLoadConfig panics when encountering errors.
func MustLoadConfig() *Config {
	cands := make([]string, 0, 2)

	if wd, err := os.Getwd(); err == nil {
		cands = append(cands, filepath.Join(wd, "config.test.yaml"))
	}
	if exe, err := os.Executable(); err == nil {
		cands = append(cands, filepath.Join(filepath.Dir(exe), "config.test.yaml"))
	}

	name := "config"
	for _, cands := range cands {
		if _, err := os.Stat(cands); err == nil {
			name = "config.test"
			break
		}
	}

	if name == "config.test" {
		slog.Info("using 'config.test.yaml' as profile")
	} else {
		slog.Info("using 'config.yaml' as profile")
	}

	cfg := confx.MustLoad[Config]("daizynight",
		confx.WithFileName(name),
		confx.WithFileType("yaml"),
	)
	return cfg

	/*if _, err := os.Stat(path); err == nil {
		slog.Info("using 'config.test.yaml' as profile")
		cfg := confx.MustLoad[Config]("daizynight",
			confx.WithFileName("config.test"),
			confx.WithFileType("yaml"),
		)
		return cfg
	} else {
		slog.Info("using 'config.yaml' as profile")
		cfg := confx.MustLoad[Config]("daizynight",
			confx.WithFileName("config"),
			confx.WithFileType("yaml"),
		)
		return cfg
	}*/
}
