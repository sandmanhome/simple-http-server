package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Config struct {
	Port int `json: "port"`
}

var config *Config
var configLock *sync.RWMutex

func LoadConfig(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	newConfig := new(Config)
	err = json.Unmarshal(data, newConfig)
	if err != nil {
		return err
	}

	configLock.Lock()
	config = newConfig
	configLock.Unlock()
	return nil
}

func GetConfig() (*Config, error) {
	configLock.RLock()
	defer configLock.RUnlock()
	if config == nil {
		return nil, fmt.Errorf("config invalid")
	}

	copyConfig := *config
	return &copyConfig, nil
}

func init() {
	configLock = new(sync.RWMutex)
}
