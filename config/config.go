package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL    string `mapstructure:"DATABASE_URL"`
	DatabaseSchema string `mapstructure:"DATABASE_SCHEMA"`
}

func LoadConfig(path string, filename string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
