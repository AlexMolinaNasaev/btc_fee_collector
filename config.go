package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
)

var (
	isLoaded   int32
	loadMux    = &sync.RWMutex{}
	globalConf *GlobalConfig
)

const (
	SatoshiInBTC   float64 = 100000000
	ByteInKilobyte float64 = 1024
)

type GlobalConfig struct {
	Node NodeFeeConfig `mapstructure:"node"`
	API  APIFeeConfig  `mapstructure:"api"`
}

type NodeFeeConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type APIFeeConfig struct {
	Url string `mapstructure:"url"`
}

// loadConfig read config from env file. if WORK_MODE env is empty or "test",
// envs will be set to default hardocoded value.
func loadConfig() error {
	loadMux.Lock()
	defer loadMux.Unlock()

	if atomic.AddInt32(&isLoaded, 1) != 1 {
		return fmt.Errorf("configs already loaded")
	}

	confPath, _ := os.Getwd()
	viper.AddConfigPath(confPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("%s: %s", "cannot read config file", err.Error())
	}

	if err := viper.Unmarshal(&globalConf); err != nil {
		return err
	}
	return nil
}

// GetConfig returns global config. If config is not already loaded will call loadConfig
func GetConfig() (*GlobalConfig, error) {
	if globalConf == nil {
		err := loadConfig()
		if err != nil {
			return nil, err
		}
	}

	return globalConf, nil
}
