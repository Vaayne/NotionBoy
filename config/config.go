package config

import "os"

type Config struct {
	BearerToken string
	DatabaseID  string
	Wechat      Wechat
}

type Wechat struct {
	AppID          string `yaml:"appID"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
}

var globalConfig = &Config{}

func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = &Config{}
	}
	return globalConfig
}

func LoadConfig(c *Config) {
	c.DatabaseID = os.Getenv("DATABASE_ID")
	c.BearerToken = os.Getenv("BEARER_TOKEN")
}
