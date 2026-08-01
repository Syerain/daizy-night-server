package confx

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func MustLoad[T any](appName string, opts ...Option) *T {
	cfg, err := Load[T](appName, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// Load loads and validates configuration
func Load[T any](appName string, opts ...Option) (*T, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	var cfg T

	defaults.SetDefaults(&cfg)

	v := viper.New()
	v.SetConfigName(o.fileName)
	v.SetConfigType(o.fileType)
	for _, path := range o.searchPaths {
		v.AddConfigPath(path)
	}

	// example: appName="myapp", field="db.host" -> "MYAPP_DB_HOST"
	v.SetEnvPrefix(appName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// ignore ConfigFileNotFoundError
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	// bind to cfg
	if err := v.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "mapstructure"
		c.ErrorUnused = true // 配置文件多余字段报错
	}); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	val := NewCustomValidator()
	if err := val.Validate(&cfg); err != nil {
		return nil, err // 直接返回 validator 的友好错误信息
	}

	return &cfg, nil
}
