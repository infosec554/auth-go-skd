package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      App      `yaml:"app"`
	HTTP     HTTP     `yaml:"http"`
	Log      Log      `yaml:"log"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
	Limiter  Limiter  `yaml:"limiter"`
	OAuth    OAuth    `yaml:"oauth"`
}

type OAuth struct {
	Google Google `yaml:"google"`
}

type Google struct {
	ClientID     string `yaml:"client_id" env:"GOOGLE_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `yaml:"redirect_url" env:"GOOGLE_REDIRECT_URL"`
}

type App struct {
	Name    string `yaml:"name" env:"APP_NAME" env-default:"auth-service"`
	Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
}

type HTTP struct {
	Port         string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-default:"auth_db"`
	SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSL_MODE" env-default:"disable"`
	PoolSize int    `yaml:"pool_size" env:"POSTGRES_POOL_SIZE" env-default:"10"`
}

type Redis struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

type Limiter struct {
	RPS   int           `yaml:"rps" env:"LIMITER_RPS" env-default:"10"`
	Burst int           `yaml:"burst" env:"LIMITER_BURST" env-default:"20"`
	TTL   time.Duration `yaml:"ttl" env:"LIMITER_TTL" env-default:"1s"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("config/config.yaml", &cfg); err != nil {
		// Fallback to reading environment variables if config file doesn't exist
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	return &cfg, nil
}
