package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTl    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, errors.New("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadByPath(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// fetchConfigPath fetches config path from environment variable or command flag
// Priority: flag > env > default = empty string
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
