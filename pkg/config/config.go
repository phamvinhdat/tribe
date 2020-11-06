package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) {
	logrus.Infof("load config from file: %s", configFile)
	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithField("error", err).Panic("failed to read in config")
	}
}
