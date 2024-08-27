package http

import (
	"net/http"
	"strconv"

	"github.com/pavr1/online_payment_platform/bank/config"
	"github.com/pavr1/online_payment_platform/bank/handlers/repo"
	"github.com/pavr1/online_payment_platform/bank/models"
	log "github.com/sirupsen/logrus"
)

type HttpHandler struct {
	log    *log.Logger
	config *config.Config
	repo   repo.IRepoHandler
}

func NewHttpHandler(log *log.Logger, config *config.Config, repoHandler repo.IRepoHandler) *HttpHandler {
	return &HttpHandler{
		log:    log,
		config: config,
		repo:   repoHandler,
	}
}

func (h *HttpHandler) Transfer() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		fromCardInfo := &models.Card{
			CardNumber: cardNumber,
			HolderName: holderName,
			ExpDate:    expDate,
			CVV:        cvv,
		}

		isValid, err := h.repo.Transfer(fromCardInfo, targetAccountNumber, float64Amount, "Payment Process")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Card is not valid"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment processed successfully"))
	}
}
