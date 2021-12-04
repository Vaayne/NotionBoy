package config

import (
	"path"

	"github.com/argoproj/pkg/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Wechat    Wechat    `mapstructure:"WECHAT"`
	Service   Service   `mapstructure:"SERVICE"`
	Databases Databases `mapstructure:"DATABASES"`
}

type Service struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

type Notion struct {
	BearerToken string `mapstructure:"BEARER_TOKEN"`
	DatabaseID  string `mapstructure:"DATABASE_ID"`
}

type Databases struct {
	Sqlite Sqlite        `mapstructure:"SQLITE"`
	MySQL  DatabaseMySQL `mapstructure:"MYSQL"`
}

type DatabaseMySQL struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Pass     string `mapstructure:"PASS"`
	Database string `mapstructure:"DATABASE"`
}

type Sqlite struct {
	File string `mapstructure:"FILE"`
}

type Wechat struct {
	AppID          string `mapstructure:"APP_ID"`
	AppSecret      string `mapstructure:"APP_SECRET"`
	Token          string `mapstructure:"TOKEN"`
	EncodingAESKey string `mapstructure:"ENCODING_AES_KEY"`
}

var globalConfig = &Config{}

func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = &Config{}
	}
	return globalConfig
}

func LoadConfig(c *Config) {
	rootPath := "."
	viper.AddConfigPath(rootPath)
	viper.SetConfigName("settings")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
	}

	if file.Exists(path.Join(rootPath, "settings_local.yaml")) {
		viper.SetConfigName("settings_local")
		err = viper.MergeInConfig()
		if err != nil {
			logrus.Error(err)
		}
	}

	viper.AutomaticEnv()
	err = viper.Unmarshal(globalConfig)
	if err != nil {
		logrus.Error(err)
	}
}
