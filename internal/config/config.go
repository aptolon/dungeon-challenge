package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}

	cfgFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(cfgFile, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
