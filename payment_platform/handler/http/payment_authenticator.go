package http

import (
	"io"
	"net/http"

	"github.com/pavr1/online_payment_platform/payment_platform/config"
	log "github.com/sirupsen/logrus"
)

type IPaymentAuthenticator interface {
	IsValidToken(token string) (int, string, error)
}
type PaymentAuthenticator struct {
	log    *log.Logger
	config *config.Config
}

func NewBankAuthenticator(log *log.Logger, config *config.Config) IPaymentAuthenticator {
	return &PaymentAuthenticator{
		log:    log,
		config: config,
	}
}

func (a *PaymentAuthenticator) IsValidToken(token string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, a.config.Auth.Path, nil)
	if err != nil {
		a.log.WithError(err).Error("Failed to create request")

		return http.StatusInternalServerError, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Entity-Key", a.config.Bank.EntityKey)

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
