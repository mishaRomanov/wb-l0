package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"pass"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	Db       string `env:"POSTGRES_DB" envDefault:"postgres"`
}

func InitConfig() (Config, error) {
	var config Config
	//here we parse all environmental variables into struct object
	err := env.Parse(&config)
	//in case it returns an error
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
