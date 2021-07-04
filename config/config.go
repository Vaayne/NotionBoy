package config

import "os"

type Config struct {
	BearerToken string
	DatabaseID  string
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
