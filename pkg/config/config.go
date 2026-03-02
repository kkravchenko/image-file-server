package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

func Config[T any](envFile string) (T, error) {
	viper.SetConfigFile(envFile)

	viper.AutomaticEnv()

	var config T

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file: ", err.Error())
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Error parsing config file: ", err.Error())
		return config, err
	}

	return config, nil
}
