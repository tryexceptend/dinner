// Пакет для работы с конфигами
package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Тип окружения
const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

// Структура конфига с привязкой к структуре из файла
type Config struct {
	Env         string `yaml:"env" env-required:"local"`
	StoragePath string `yaml:"storage_path" env-required:"./data"`
	Timeout     int    `yaml:"timeout" env-default:30`
}

// MustLoadConfig загружает конфиг из файла в структуру Config
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

// fetchConfigPath берет путь к файлу конфига из аргументов или из переменной окружения
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("DINNER_CONFIG_PATH")
	}
	return res
}
