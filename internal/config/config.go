package config

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gcfg"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	MQTT   MQTTConfig   `yaml:"mqtt"`
	Logger Logger       `yaml:"logger"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type SQLiteConfig struct {
	Path     string `yaml:"path"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type MQTTConfig struct {
	Address string
}
type Logger struct {
	Path       string `yaml:"path"`
	File       string `yaml:"file"`
	Level      string `yaml:"level"`
	TimeFormat string `yaml:"timeFormat"`
	Stdout     bool   `yaml:"stdout"`
	Header     bool   `yaml:"header"`
}

// Load 从指定目录读取 config.yaml 并映射到结构体。
func Load(ctx context.Context, configDir string) (*Config, error) {
	adapter, err := gcfg.NewAdapterFile(configDir)
	if err != nil {
		return nil, fmt.Errorf("create config adapter: %w", err)
	}
	cfg := gcfg.NewWithAdapter(adapter)

	var out Config
	if err := cfg.MustGet(ctx, ".").Struct(&out); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &out, nil
}
