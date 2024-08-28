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
	Auth struct {
		Path string
	}
	MongoDB struct {
		Uri                    string
		Database               string
		Card_Collection        string
		Transaction_Collection string
		//pvillalobos add this to a secret later
		Username string
		Password string
		//RolName  string
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

	authPath := os.Getenv("AUTH_PATH")
	if authPath == "" {
		log.Error("AUTH_PATH is not set")
		return nil, errors.New("AUTH_PATH is not set")
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		log.Error("MONGODB_URI is not set")
		return nil, errors.New("MONGODB_URI is not set")
	}

	mongodb_database := os.Getenv("MONGODB_DATABASE")
	if mongodb_database == "" {
		log.Error("MONGODB_DATABASE is not set")
		return nil, errors.New("MONGODB_DATABASE is not set")
	}

	mongodb_card_collection := os.Getenv("MONGODB_CARD_COLLECTION")
	if mongodb_card_collection == "" {
		log.Error("MONGODB_CARD_COLLECTION is not set")
		return nil, errors.New("MONGODB_CARD_COLLECTION is not set")
	}

	mongodb_transaction_collection := os.Getenv("MONGODB_TRANSACTION_COLLECTION")
	if mongodb_transaction_collection == "" {
		log.Error("MONGODB_TRANSACTION_COLLECTION is not set")
		return nil, errors.New("MONGODB_TRANSACTION_COLLECTION is not set")
	}

	mongodb_username := os.Getenv("MONGODB_USERNAME")
	if mongodb_username == "" {
		log.Error("MONGODB_USERNAME is not set")
		return nil, errors.New("MONGODB_USERNAME is not set")
	}

	mongodb_password := os.Getenv("MONGODB_PASSWORD")
	if mongodb_password == "" {
		log.Error("MONGODB_PASSWORD is not set")
		return nil, errors.New("MONGODB_PASSWORD is not set")
	}

	var config = Config{}
	config.Server.Port = portInt
	config.Auth.Path = authPath
	config.MongoDB.Uri = mongodb_uri
	config.MongoDB.Database = mongodb_database
	config.MongoDB.Card_Collection = mongodb_card_collection
	config.MongoDB.Transaction_Collection = mongodb_transaction_collection
	config.MongoDB.Username = mongodb_username
	config.MongoDB.Password = mongodb_password

	log.WithField("config", config).Info("Loaded environment variables")

	return &config, nil
}
