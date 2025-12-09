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
		SSLMode  string
	}
	JWT struct {
		Secret    string
		ExpiresIn int
	}
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	return cfg
}
