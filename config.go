package main

import (
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/spf13/viper"
)

func LoadConfig() error {
	viper.AddConfigPath("src/github.com/bmotto/go_shortlink")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}
