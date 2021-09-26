package config

import (
	"encoding/json"
	"os"

	"github.com/morzik45/test-go/pkg/repository"
)

type ServerHttpConfig struct {
	Port string
}

type AppConfig struct {
	Db   repository.Config
	Http ServerHttpConfig
}

func InitConfig(path string) (*AppConfig, error) {
	var config AppConfig
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
