package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port      int    `mapstructure:"port"`
		SecretKey string `mapstructure:"secret"`
	} `mapstructure:"server"`
}

func NewConfig(log *log.Logger) (*Config, error) {
	// Unmarshal the configuration into a struct
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		log.Error("AUTH_PORT is not set")
		return nil, errors.New("AUTH_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithField("error", err).Error("Failed to convert port to int")
		return nil, err
	}

	secretKey := os.Getenv("AUTH_SECRET_KEY")
	var config = Config{}
	config.Server.Port = portInt
	config.Server.SecretKey = secretKey

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
