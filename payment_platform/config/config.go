package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port int
	}
	Bank struct {
		Host      string
		EntityKey string
	}
	Auth struct {
		Path string
	}
}

func NewConfig() (*Config, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Error("SERVER_PORT is not set")
		return nil, errors.New("SERVER_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithField("error", err).Error("Failed to convert port to int")
		return nil, err
	}

	bankHost := os.Getenv("BANK_HOST")
	if bankHost == "" {
		log.Error("BANK_HOST is not set")
		return nil, errors.New("BANK_HOST is not set")
	}

	bankEntityKey := os.Getenv("BANK_ENTITY_KEY")
	if bankEntityKey == "" {
		log.Error("BANK_ENTITY_KEY is not set")
		return nil, errors.New("BANK_ENTITY_KEY is not set")
	}

	authPath := os.Getenv("AUTH_PATH")
	if authPath == "" {
		log.Error("AUTH_PATH is not set")
		return nil, errors.New("AUTH_PATH is not set")
	}

	var config = Config{}
	config.Server.Port = portInt
	config.Bank.Host = bankHost
	config.Bank.EntityKey = bankEntityKey
	config.Auth.Path = authPath

	log.WithField("config", config).Info("Loaded environment variables")

	return &config, nil
}
