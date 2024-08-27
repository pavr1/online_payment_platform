package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Auth struct {
		Path string `mapstructure:"path"`
		Host string `mapstructure:"host"`
	} `mapstructure:"auth"`
	MongoDB struct {
		Uri        string `mapstructure:"uri"`
		Database   string `mapstructure:"database"`
		Collection string `mapstructure:"collection"`
		//pvillalobos add this to a secret later
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		RolName  string `mapstructure:"role"`
	} `mapstructure:"mongodb"`
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

	var config = Config{}
	config.Server.Port = portInt

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
