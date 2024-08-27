package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pavr1/online_payment_platform/bank/config"
	_http "github.com/pavr1/online_payment_platform/bank/handlers/http"
	"github.com/pavr1/online_payment_platform/bank/handlers/repo"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()

	log := setupLogger()
	config, err := config.NewConfig()
	if err != nil {
		return
	}

	repoHandler, err := repo.NewRepoHandler(log, config)
	if err != nil {
		return
	}
	httpHandler := _http.NewHttpHandler(log, config, repoHandler)

	router.HandleFunc("/transfer", httpHandler.Transfer())

	log.WithField("port", config.Server.Port).Info("Listening to Server...")
	// Start the HTTP server
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router))
}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)

	// Set the output to stdout
	logger.SetOutput(os.Stdout)

	return logger
}
