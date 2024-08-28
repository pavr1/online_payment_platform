package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port                     int    `mapstructure:"port"`
		BankSecretKey            string `mapstructure:"bank_secret_key"`
		PaymentPlatformSecretKey string `mapstructure:"payment_platform_secret_key"`
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

	bankSecretKey := os.Getenv("AUTH_BANK_SECRET_KEY")
	if bankSecretKey == "" {
		log.Error("AUTH_BANK_SECRET_KEY is not set")
		return nil, errors.New("AUTH_BANK_SECRET_KEY is not set")
	}

	PaymentPlatformSecretKey := os.Getenv("AUTH_PAYMENT_PLATFORM_SECRET_KEY")
	if PaymentPlatformSecretKey == "" {
		log.Error("AUTH_PAYMENT_PLATFORM_SECRET_KEY is not set")
		return nil, errors.New("AUTH_PAYMENT_PLATFORM_SECRET_KEY is not set")
	}

	var config = Config{}
	config.Server.Port = portInt
	config.Server.BankSecretKey = bankSecretKey
	config.Server.PaymentPlatformSecretKey = PaymentPlatformSecretKey

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
