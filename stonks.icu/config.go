package main

import "github.com/spf13/viper"

func configure() {
	viper.SetDefault("port", 8080)
	viper.AutomaticEnv()
}
