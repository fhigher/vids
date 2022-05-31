package config

import (
	"github.com/jinzhu/configor"
)

const (
	ConfigPathEnv = "VIDS_SERVER_CONFIG"
	Prod          = "production"
	Dev           = "development"
	Test          = "test"
)

type Config struct {
	ServerName string `default:"Vids Server"` // server name
	ServerListen     string
	ServerPort       string `default:"9000"`
}

func InitConfig(configPath string, c *configor.Config) (config Config, err error) {
	err = configor.New(c).Load(&config, configPath)
	return
}
