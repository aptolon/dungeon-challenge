package config

import (
	"encoding/json"
	"os"
	"time"
)

type rawConfig struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

type Config struct {
	Floors   int
	Monsters int
	OpenAt   time.Time
	Duration int
}

func LoadConfig(path string) (Config, error) {
	var raw rawConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{}, err
	}

	openAt, err := time.Parse("15:04:05", raw.OpenAt)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Floors:   raw.Floors,
		Monsters: raw.Monsters,
		OpenAt:   openAt,
		Duration: raw.Duration,
	}, nil
}
