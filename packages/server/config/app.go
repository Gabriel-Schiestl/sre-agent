package config

type AppConfig struct {
	Port int `env:"APP_PORT" default:"8080"`
}

func LoadApp() (*AppConfig, error) {
	cfg := &AppConfig{}
	err := LoadConfig(cfg)
	return cfg, err
}

func LoadConfig(cfg any) error {
	return nil
}