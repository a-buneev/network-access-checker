package models

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	sysConfig
	checkerConfig
}

type checkerConfig struct {
	ResourceList           []Resource `json:"resourceList"`
	CheckPeriodSeconds     int        `json:"checkPeriodSeconds"`
	CheckConnectionTimeout int        `json:"checkConnectionTimeout"`
}

type sysConfig struct {
	mu       sync.Mutex
	filePath string
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
	config.filePath = path
	return &config
}

func (conf *Config) Reload(period int) {
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	for {
		<-ticker.C
		file, err := os.ReadFile(conf.filePath)
		if err != nil {
			log.Printf("Reload config failed (read file): %v", err.Error())
		}
		var tempConfig checkerConfig
		err = json.Unmarshal(file, &tempConfig)
		if err != nil {
			log.Printf("Reload config failed (unmarshall): %v", err.Error())
		}
		conf.mu.Lock()
		conf.checkerConfig = tempConfig
		conf.mu.Unlock()
	}
}
