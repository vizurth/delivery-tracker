package config

import (
	"delivery-tracker/common/kafka"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type NotificationsConfig struct {
	Port int `yaml:"port" env:"PORT"`
}

type Config struct {
	Notification NotificationsConfig `yaml:"auth" env-previx:"AUTH_"`
	Kafka        kafka.KafkaConfig   `yaml:"kafka" env-previx:"KAFKA"`
}

func New() (Config, error) {
	var config Config
	// docker workdir app/
	// local workdir delivery-tracker/notification
	if err := cleanenv.ReadConfig("configs/config.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading config: %w", err)
		}
	}

	return config, nil
}
