package config

import (
	"delivery-tracker/common/postgres"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type AuthConfig struct {
	Port          int    `yaml:"port" env:"PORT"`
	MigrationPath string `yaml:"migration-path" env:"MIGRATION_PATH"`
	SecretKey     string `yaml:"secret-key" env:"SECRET_KEY" env-default:"your-secret-key"`
}

type Config struct {
	Postgres postgres.Config `yaml:"postgres" env-prefix:"POSTGRES_"`
	Auth     AuthConfig      `yaml:"auth" env-previx:"AUTH_"`
}

func New() (Config, error) {
	var config Config
	// docker workdir app/
	// local workdir delivery-tracker/auth
	if err := cleanenv.ReadConfig("configs/config.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading config: %w", err)
		}
	}

	return config, nil
}
