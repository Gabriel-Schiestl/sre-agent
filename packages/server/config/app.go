package config

import (
	"os"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	Port            int    `env:"APP_PORT"          envDefault:"8080"`
	UploadsDir      string `env:"UPLOADS_DIR"       envDefault:"./data/uploads"`
	FrontendURL     string `env:"FRONTEND_URL"      envDefault:"http://localhost:3000"`
	AnthropicAPIKey string `env:"ANTHROPIC_API_KEY"`
}

func LoadApp() (*AppConfig, error) {
	cfg := &AppConfig{}
	err := LoadConfig(cfg)
	return cfg, err
}

func LoadConfig(cfg any) error {
	return env.Parse(cfg)
}

func EnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
