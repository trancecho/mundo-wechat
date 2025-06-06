package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

var GlobalConfig *Config

func Init() {
	GlobalConfig = &Config{
		viper.New(),
	}
	GlobalConfig.SetConfigName("application")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath(".")
	GlobalConfig.AddConfigPath("./config")

	err := GlobalConfig.ReadInConfig()
	if err != nil {
		logrus.WithField("config", "GlobalConfig").WithError(err).Panicf("unable to read global config")
	}

	GlobalConfig.WatchConfig()
	GlobalConfig.OnConfigChange(func(e fsnotify.Event) {
		err := GlobalConfig.ReadInConfig()
		if err == nil {
			logrus.WithField("config", "GlobalConfig").Info("config updated")
		} else {
			logrus.WithField("config", "GlobalConfig").WithError(err).Error("failed to update config")
		}
	})
}
