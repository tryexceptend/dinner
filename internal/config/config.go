package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env         string `yaml:"env" env-required:"local"`
	StoragePath string `yaml:"storage_path" env-required:"./data"`
	Timeout     int    `yaml:"timeout" env-default:30`
}

func MustLoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

// берет путь к файлу конфига из аргументов
// или из переменной окружения
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("DINNER_CONFIG_PATH")
	}
	return res
}
