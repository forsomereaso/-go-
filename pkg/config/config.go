package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	App   App   `mapstructure:"app"`
	Log   Log   `mapstructure:"log"`
	MySQL MySQL `mapstructure:"mysql"`
	Redis Redis `mapstructure:"redis"`
}

type App struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type MySQL struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Init() error {
	var env string
	flag.StringVar(&env, "env", "dev", "运行环境: dev/test/prod")
	flag.Parse()

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config err: %v", err)
	}

	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("unmarshal config err: %v", err)
	}

	return nil
}
