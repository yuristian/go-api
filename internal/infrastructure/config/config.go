package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	DB struct {
		Type     string
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret    string
		ExpiresIn int
	}
}

func LoadConfig(path string) *Config {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	return &config
}
