package config

type DBConfig struct {
	Driver   string `env:"DB_DRIVER" default:"postgres"`
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER" default:"postgres"`
	Password string `env:"DB_PASSWORD" default:"password"`
	Name     string `env:"DB_NAME" default:"sre_agent"`
}

func LoadDB() (*DBConfig, error) {
	cfg := &DBConfig{}
	err := LoadConfig(cfg)
	return cfg, err
}