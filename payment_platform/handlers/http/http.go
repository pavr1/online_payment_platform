package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavr1/online_payment_platform/payment_platform/config"
	"github.com/pavr1/online_payment_platform/payment_platform/handlers/http/providers"
	log "github.com/sirupsen/logrus"
)

type HttpHandler struct {
	log           *log.Logger
	config        *config.Config
	tokenProvider providers.ITokenProvider
	bankProvider  providers.IBankProvider
	client        *http.Client
}

func NewHttpHandler(log *log.Logger, config *config.Config, tokenProvider providers.ITokenProvider, bankProvider providers.IBankProvider, client *http.Client) *HttpHandler {
	return &HttpHandler{
		log:           log,
		config:        config,
		tokenProvider: tokenProvider,
		client:        client,
		bankProvider:  bankProvider,
	}
}

func (h *HttpHandler) ProcessPurchase() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.log.Error("Method not allowed")

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		h.log.Info("Handling purchase request...")

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			h.log.Error("Token is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token is required"))
			return
		}

		statusCode, body, err := h.tokenProvider.IsValidToken(token)
		if err != nil {
			h.log.WithError(err).Error("Failed to validate token")

			w.WriteHeader(statusCode)
			w.Write([]byte("Failed to validate token"))
			return
		}

		if statusCode != http.StatusOK {
			h.log.WithField("StatusCode", statusCode).WithField("Body", body).Error("Failed to validate token")

			w.WriteHeader(statusCode)
			w.Write([]byte(body))
			return
		}

		cardNumber := r.Header.Get("card_number")
		if cardNumber == "" {
			h.log.Error("Card number is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Card number is required"))
			return
		}

		holderName := r.Header.Get("holder_name")
		if holderName == "" {
			h.log.Error("Holder name is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Holder name is required"))
			return
		}
		expDate := r.Header.Get("exp_date")
		if expDate == "" {
			h.log.Error("Expiration date is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Expiration date is required"))
			return
		}
		cvv := r.Header.Get("cvv")
		if cvv == "" {
			h.log.Error("CVV is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("CVV is required"))
			return
		}

		targetAccountNumber := r.Header.Get("target_account_number")
		if targetAccountNumber == "" {
			h.log.Error("Target account number is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Target account number is required"))
			return
		}

		amount := r.Header.Get("amount")
		if amount == "" {
			h.log.Error("Amount is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Amount is required"))
			return
		}

		float64Amount, err := strconv.ParseFloat(amount, 32)
		if err != nil {
			h.log.WithError(err).Error("Failed to convert amount to float32")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to convert amount to float32"))
			return
		}

		//todo: this token should be used only for validating payment_platform access, a brand new one should be created for bank authorization
		statusCode, body, err = h.bankProvider.ProcessPayment(token, cardNumber, holderName, expDate, cvv, targetAccountNumber, float64Amount)
		if err != nil {
			h.log.WithError(err).Error("Failed to process payment with the bank")

			w.WriteHeader(statusCode)
			w.Write([]byte("Failed to process payment with the bank"))
			return
		}

		if statusCode != http.StatusOK {
			h.log.WithField("StatusCode", statusCode).WithField("Body", body).Error("Processed payment with the bank failed")

			w.WriteHeader(statusCode)
			w.Write([]byte(body))
			return
		}

		h.log.Info("Payment processed successfully.")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment processed successfully.\n" + body))
	}
}

func (h *HttpHandler) GetTransactionHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.log.Error("Method not allowed")

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		h.log.Info("Handling retrieve request...")

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			h.log.Error("Token is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token is required"))
			return
		}

		statusCode, body, err := h.tokenProvider.IsValidToken(token)
		if err != nil {
			h.log.WithError(err).Error("Failed to validate token")

			w.WriteHeader(statusCode)
			w.Write([]byte("Failed to validate token"))
			return
		}

		if statusCode != http.StatusOK {
			h.log.WithField("StatusCode", statusCode).WithField("Body", body).Error("Failed to validate token")

			w.WriteHeader(statusCode)
			w.Write([]byte(body))
			return
		}

		accountNumber := r.Header.Get("account_number")
		if accountNumber == "" {
			h.log.Error("account_number is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("account_number is required"))
			return
		}

		transactions, err := h.bankProvider.GetHistory(token, accountNumber)
		if err != nil {
			h.log.WithError(err).Error("Failed getting transaction los from the bank")

			w.WriteHeader(statusCode)
			w.Write([]byte("Failed getting transaction los from the bank"))
			return
		}

		if statusCode != http.StatusOK {
			h.log.WithField("StatusCode", statusCode).WithField("Body", body).Error("Transaction logs fetch from the bank failed")

			w.WriteHeader(statusCode)
			w.Write([]byte(body))
			return
		}

		transactionLogs, err := json.Marshal(transactions)
		if err != nil {
			h.log.WithError(err).Error("Failed to marshal transaction logs")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal transaction logs"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(transactionLogs)
	}
}
