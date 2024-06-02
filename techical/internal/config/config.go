package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	FastForexAPIKey string `yaml:"forex_api_key" env:"FOREX_API_KEY"`
	ListenPort      string `yaml:"listen_port" env:"LISTEN_PORT"`
	DatabaseHost    string `yaml:"database_host" env:"DATABASE_HOST"`
	DatabasePort    string `yaml:"database_port" env:"DATABASE_PORT"`
	DatabaseUser    string `yaml:"database_user" env:"DATABASE_USER"`
	DatabasePass    string `yaml:"database_password" env:"DATABASE_PASS"`
	DatabaseName    string `yaml:"database_name" env:"DATABASE_NAME"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
