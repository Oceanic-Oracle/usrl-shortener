package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `env:"ENV"`
	Port    string `env:"PORT"`
	Storage StorageCfg
}

type StorageCfg struct {
	StorageType string `env:"STORAGE_TYPE"`
	Port        string `env:"PORT"`
	User        string `env:"POSTGRES_USER"`
	Password    string `env:"POSTGRES_PASSWORD"`
	DataBase    string `env:"POSTGRES_DB"`
}

func MustLoad(path string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
		return nil
	}

	return &cfg
}
