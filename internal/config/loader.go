package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Mongo struct {
		URI        string `mapstructure:"uri"`
		Database   string `mapstructure:"database"`
		Collection string `mapstructure:"collection"`
	} `mapstructure:"mongo"`
	Kubernetes struct {
		Namespace string `mapstructure:"namespace"`
	} `mapstructure:"kubernetes"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".") // Look in the current directory
	viper.AutomaticEnv()     // Automatically override with environment variables

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config file: %v. Proceeding with defaults/env vars.", err)
	}

	// Unmarshal config into a struct
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
