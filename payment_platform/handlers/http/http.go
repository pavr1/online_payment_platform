package http

import (
	"net/http"
	"strings"

	"github.com/pavr1/online_payment_platform/payment_platform/config"
	"github.com/pavr1/online_payment_platform/payment_platform/handlers/http/providers"
	log "github.com/sirupsen/logrus"
)

type HttpHandler struct {
	log           *log.Logger
	config        *config.Config
	tokenProvider providers.ITokenProvider
	client        *http.Client
}

func NewHttpHandler(log *log.Logger, config *config.Config, tokenProvider providers.ITokenProvider, client *http.Client) *HttpHandler {
	return &HttpHandler{
		log:           log,
		config:        config,
		tokenProvider: tokenProvider,
		client:        client,
	}
}

func (h *HttpHandler) ProcessPurchase() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			h.log.Error("Token is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token is required"))
			return
		}

		statusCode, body, err := h.tokenProvider.IsValidToken(token)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte("Failed to validate token"))
			return
		}

		if statusCode != http.StatusOK {
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

		req, err := http.NewRequest(http.MethodGet, h.config.Bank.Path, nil)
		if err != nil {
			h.log.WithError(err).Error("Failed to create request")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create request"))
			return
		}

		//todo: encrypt all this data
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("card_number", cardNumber)
		req.Header.Set("holder_name", holderName)
		req.Header.Set("exp_date", expDate)
		req.Header.Set("cvv", cvv)
		req.Header.Set("target_account_number", targetAccountNumber)
		req.Header.Set("amount", amount)
		req.Header.Set("X-Entity-Key", "YmFuay1zZWNyZXQta2V5LWF1dGhlbnRpY2F0aW9u")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			h.log.WithField("Path", h.config.Auth.Path).WithError(err).Error("Failed to send request to the bank")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to comunicate wuth the bank"))
			return
		}

		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment processed successfully"))
	}
}
