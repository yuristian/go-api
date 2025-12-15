package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
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

	cfg.validate()

	return cfg
}

func (c *Config) validate() {
	if c.JWT.Secret == "" {
		log.Fatal("JWT secret is required")
	}
	if c.JWT.ExpiresIn <= 0 {
		log.Fatal("JWT expires_in must be > 0")
	}
	if c.Server.Port == 0 {
		log.Fatal("Server port is required")
	}
}
