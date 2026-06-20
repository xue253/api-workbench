package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

var AppConfig Config

func Load(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&AppConfig)
}
