package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"INFO"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("LOG", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
