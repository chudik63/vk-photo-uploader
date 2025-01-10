package infrastructure

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig("./configs/.env", &cfg)

	if cfg == (Config{}) {
		return nil, errors.New("конфиг пустой")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
