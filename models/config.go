package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type Config struct {
	SysConfig
	CheckerConfig
}

type CheckerConfig struct {
	ResourceList           []Resource `json:"resourceList"`
	CheckPeriodSeconds     int        `json:"checkPeriodSeconds"`
	CheckConnectionTimeout int        `json:"checkConnectionTimeout"`
}

type SysConfig struct {
	mu       sync.Mutex
	filePath string
}

func NewConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
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

func (conf *Config) Reload(period int, reload chan bool) {
	fmt.Println("Start config autoreloader")
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	for {
		<-ticker.C
		//fmt.Println("Check config")
		file, err := ioutil.ReadFile(conf.filePath)
		if err != nil {
			log.Printf("Reload config failed (read file): %v", err.Error())
		}
		var tempConfig CheckerConfig
		err = json.Unmarshal(file, &tempConfig)
		if err != nil {
			log.Printf("Reload config failed (unmarshall): %v", err.Error())
		}
		if fmt.Sprintf("%v", conf.CheckerConfig) == fmt.Sprintf("%v", tempConfig) {
			//fmt.Println("Check config - no changes")
			continue
		}
		fmt.Println("Config has changed!")
		conf.mu.Lock()
		conf.CheckerConfig = tempConfig
		reload <- true
		conf.mu.Unlock()
	}
}
