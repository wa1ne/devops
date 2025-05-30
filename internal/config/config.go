package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"dev"`
	Probability4xx float32    `yaml:"probability4xx" env-default:"0.5"`
	Probability5xx float32    `yaml:"probability5xx" env-default:"1.0"`
	Server         HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := "./config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error(
			"конфиг файл не найден",
			slog.String("path", configPath),
		)
		os.Exit(1)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		slog.Error(
			"ошибка при чтении конфига",
			"error", err,
		)
		os.Exit(1)
	}

	return cfg
}
