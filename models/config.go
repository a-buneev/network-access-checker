package models

import (
	"encoding/json"
	"os"
)

type Config struct {
	ResourceList []Resource `json:"resourceList"`
	MetricsPort  string     `json:"metricsPort"`
}

func NewConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
