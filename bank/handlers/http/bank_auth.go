package http

import (
	"io"
	"net/http"

	"github.com/pavr1/online_payment_platform/bank/config"
	log "github.com/sirupsen/logrus"
)

type IBankAuthenticator interface {
	IsValidToken(token string) (int, string, error)
}
type BankAuthenticator struct {
	log    *log.Logger
	config *config.Config
}

func NewBankAuthenticator(log *log.Logger, config *config.Config) IBankAuthenticator {
	return &BankAuthenticator{
		log:    log,
		config: config,
	}
}

func (a *BankAuthenticator) IsValidToken(token string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, a.config.Auth.Path, nil)
	if err != nil {
		a.log.WithError(err).Error("Failed to create request")

		return http.StatusInternalServerError, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.log.WithField("Path", a.config.Auth.Path).WithError(err).Error("Failed to send request")

		return http.StatusInternalServerError, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.WithError(err).Error("Failed to read response")

		return http.StatusInternalServerError, "", err
	}

	return resp.StatusCode, string(body), nil
}
