package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pavr1/online_payment_platform/payment_platform/config"
	"github.com/pavr1/online_payment_platform/payment_platform/models"
	log "github.com/sirupsen/logrus"
)

type IBankProvider interface {
	ProcessPayment(token, cardNumber, holderName, expDate, cvv, targetAccountNumber string, amount float64) (int, string, error)
	ProcessRefund(token, referenceNumber string) (int, string, error)
	GetHistory(token, accountNumber string) ([]*models.Transaction, error)
	CreateBankToken(userName string) (string, error)
}
type BankProvider struct {
	log    *log.Logger
	config *config.Config
}

func NewBankProvider(log *log.Logger, config *config.Config) IBankProvider {
	return &BankProvider{
		log:    log,
		config: config,
	}
}

func (b *BankProvider) ProcessPayment(token, cardNumber, holderName, expDate, cvv, targetAccountNumber string, amount float64) (int, string, error) {
	endpoint := fmt.Sprintf("%s/transfer", b.config.Bank.Host)

	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to create request")

		return http.StatusInternalServerError, "", err
	}

	//todo: encrypt all this data
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("card_number", cardNumber)
	req.Header.Set("holder_name", holderName)
	req.Header.Set("exp_date", expDate)
	req.Header.Set("cvv", cvv)
	req.Header.Set("target_account_number", targetAccountNumber)
	req.Header.Set("amount", strconv.FormatFloat(amount, 'f', 6, 64))
	req.Header.Set("X-Entity-Key", "YmFuay1zZWNyZXQta2V5LWF1dGhlbnRpY2F0aW9u")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to send request")

		return http.StatusInternalServerError, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.log.WithError(err).Error("Failed to read response")

		return http.StatusInternalServerError, "", err
	}

	return resp.StatusCode, string(body), nil
}

func (b *BankProvider) GetHistory(token, accountNumber string) ([]*models.Transaction, error) {
	endpoint := fmt.Sprintf("%s/history", b.config.Bank.Host)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to create request")

		return []*models.Transaction{}, err
	}

	//todo: encrypt all this data
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("account_number", accountNumber)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to send request")

		return []*models.Transaction{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.log.WithError(err).Error("Failed to read response")

		return []*models.Transaction{}, err
	}

	transactions := []*models.Transaction{}
	err = json.Unmarshal(body, &transactions)
	if err != nil {
		b.log.WithError(err).Error("Failed to unmarshal response")

		return []*models.Transaction{}, err
	}

	return transactions, nil
}

func (b *BankProvider) ProcessRefund(token, referenceNumber string) (int, string, error) {
	endpoint := fmt.Sprintf("%s/refund", b.config.Bank.Host)

	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to create request")

		return http.StatusInternalServerError, "", err
	}

	//todo: encrypt all this data
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("reference_number", referenceNumber)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.log.WithField("Path", endpoint).WithError(err).Error("Failed to send request")

		return http.StatusInternalServerError, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.log.WithError(err).Error("Failed to read response")

		return http.StatusInternalServerError, "", err
	}

	return resp.StatusCode, string(body), nil
}

func (b *BankProvider) CreateBankToken(userName string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, b.config.Auth.Path, nil)
	if err != nil {
		b.log.WithField("Path", b.config.Auth.Path).WithError(err).Error("Failed to create request")

		return "", err
	}

	//todo: encrypt all this data
	req.Header.Set("X-User-Name", userName)
	req.Header.Set("X-Entity-Name", "Bank")
	req.Header.Set("X-Entity-Key", b.config.Bank.EntityKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.log.WithField("Path", b.config.Auth.Path).WithError(err).Error("Failed to send request")

		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.log.WithError(err).Error("Failed to read response")

		return "", err
	}

	return string(body), nil
}
