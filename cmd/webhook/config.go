package main

import (
	"fmt"

	"github.com/cmaster11/k8s-event-watcher/pkg"
	"github.com/spf13/viper"
)

type Webhook struct {
	Url     string            `yaml:"url" validate:"required,url"`
	Headers map[string]string `yaml:"headers"`
}

type Config struct {
	pkg.Config `mapstructure:",squash"`

	Webhooks []*Webhook `mapstructure:"webhooks" yaml:"webhooks"`
}

func (c *Config) Validate() error {
	if err := pkg.Validate.Struct(c); err != nil {
		return err
	}

	return nil
}

func loadConfigEnv(configFile *string) *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	if configFile != nil {
		v.SetConfigFile(*configFile)
	}
	v.SetConfigType("yaml")
	v.SetEnvPrefix("k8sew")
	v.AutomaticEnv()
	return v
}

func parseConfig(configPath *string) (*Config, error) {
	vip := loadConfigEnv(configPath)
	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	config := new(Config)

	if err := vip.Unmarshal(config, viper.DecodeHook(pkg.StringToRegexHookFunc())); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
