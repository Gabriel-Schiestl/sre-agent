package config

type DBConfig struct {
	Driver   string `env:"DB_DRIVER"   envDefault:"postgres"`
	Host     string `env:"DB_HOST"     envDefault:"localhost"`
	Port     int    `env:"DB_PORT"     envDefault:"5432"`
	User     string `env:"DB_USER"     envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
	Name     string `env:"DB_NAME"     envDefault:"sre_agent"`
}

func LoadDB() (*DBConfig, error) {
	cfg := &DBConfig{}
	err := LoadConfig(cfg)
	return cfg, err
}
