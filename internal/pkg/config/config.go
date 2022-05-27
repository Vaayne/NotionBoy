package config

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Wechat      Wechat      `mapstructure:"WECHAT"`
	Service     Service     `mapstructure:"SERVICE"`
	Databases   Databases   `mapstructure:"DATABASES"`
	NotionOauth NotionOauth `mapstructure:"NOTION_OAUTH"`
	R2Config    R2Config    `mapstructure:"R2_CONFIG"`
}

type Service struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
	URL  string `mapstructure:"URL"`
}

type Notion struct {
	BearerToken string `mapstructure:"BEARER_TOKEN"`
	DatabaseID  string `mapstructure:"DATABASE_ID"`
}

type NotionOauth struct {
	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	RedirectURI  string `mapstructure:"REDIRECT_URI"`
	AuthURL      string `mapstructure:"AUTH_URL"`
	TokenURL     string `mapstructure:"TOKEN_URL"`
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

type R2Config struct {
	Token string `mapstructure:"TOKEN"`
	Url   string `mapstructure:"URL"`
}

var globalConfig = &Config{}

func init() {
	logrus.Debug("init config")
	LoadConfig(globalConfig)
}

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

	if fileExists(path.Join(rootPath, "settings_local.yaml")) {
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
