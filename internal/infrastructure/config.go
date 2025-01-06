package infrastructure

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port string
	}
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)

	if cfg == (Config{}) {
		return nil, errors.New("конфиг пустой")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
