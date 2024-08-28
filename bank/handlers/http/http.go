package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavr1/online_payment_platform/bank/config"
	"github.com/pavr1/online_payment_platform/bank/handlers/repo"
	"github.com/pavr1/online_payment_platform/bank/models"
	log "github.com/sirupsen/logrus"
)

type HttpHandler struct {
	log               *log.Logger
	config            *config.Config
	repo              repo.IRepoHandler
	bankAuthenticator IBankAuthenticator
}

func NewHttpHandler(log *log.Logger, config *config.Config, repoHandler repo.IRepoHandler, bankAuthenticator IBankAuthenticator) *HttpHandler {
	return &HttpHandler{
		log:               log,
		config:            config,
		repo:              repoHandler,
		bankAuthenticator: bankAuthenticator,
	}
}

func (h *HttpHandler) Transfer() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			h.log.Error("Token is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token is required"))
			return
		}

		statusCode, body, err := h.bankAuthenticator.IsValidToken(token)
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

		status, referenceNumber, err := h.repo.Transfer(fromCardInfo, targetAccountNumber, float64Amount, "Purchase Processed")
		if err != nil {
			w.WriteHeader(status)
			w.Write([]byte("failed making transfer"))
			return
		}

		if status != http.StatusOK {
			w.WriteHeader(status)
			w.Write([]byte("Card is not valid"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reference Numer: " + referenceNumber))
	}
}

func (h *HttpHandler) Fillup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cards := []*models.Card{
			{
				ID:         "12345",
				CardNumber: "4532-1143-8765-3211",
				HolderName: "Emily Chen",
				ExpDate:    "02/2027",
				CVV:        "987",
				Account: &models.Account{
					ID:            "account-001",
					AccountNumber: "1234567890",
					Amount:        1234.56,
				},
				Customer: &models.Customer{
					ID:        "customer-001",
					FirstName: "Emily",
					LastName:  "Chen",
					Email:     "emily.chen@example.com",
				},
			},
			{
				ID:         "67890",
				CardNumber: "9876-5432-1098-7654",
				HolderName: "David Lee",
				ExpDate:    "12/2024",
				CVV:        "654",
				Account: &models.Account{
					ID:            "account-002",
					AccountNumber: "9876543210",
					Amount:        987.65,
				},
				Customer: &models.Customer{
					ID:        "customer-002",
					FirstName: "David",
					LastName:  "Lee",
					Email:     "david.lee@example.com",
				},
			},
			{
				ID:         "34567",
				CardNumber: "2345-6789-0123-4567",
				HolderName: "Sophia Patel",
				ExpDate:    "06/2026",
				CVV:        "321",
				Account: &models.Account{
					ID:            "account-003",
					AccountNumber: "1111111111",
					Amount:        1111.11,
				},
				Customer: &models.Customer{
					ID:        "customer-003",
					FirstName: "Sophia",
					LastName:  "Patel",
					Email:     "sophia.patel@example.com",
				},
			},
			{
				ID:         "90123",
				CardNumber: "5678-9012-3456-7890",
				HolderName: "Michael Kim",
				ExpDate:    "03/2025",
				CVV:        "876",
				Account: &models.Account{
					ID:            "account-004",
					AccountNumber: "2222222222",
					Amount:        222.22,
				},
				Customer: &models.Customer{
					ID:        "customer-004",
					FirstName: "Michael",
					LastName:  "Kim",
					Email:     "michael.kim@example.com",
				},
			},
			{
				ID:         "45678",
				CardNumber: "8901-2345-6789-0123",
				HolderName: "Olivia Brown",
				ExpDate:    "09/2026",
				CVV:        "543",
				Account: &models.Account{
					ID:            "account-005",
					AccountNumber: "3333333333",
					Amount:        3333.33,
				},
				Customer: &models.Customer{
					ID:        "customer-005",
					FirstName: "Olivia",
					LastName:  "Brown",
					Email:     "olivia.brown@example.com",
				},
			},
			{
				ID:         "23456",
				CardNumber: "1234-5678-9012-3456",
				HolderName: "William White",
				ExpDate:    "01/2027",
				CVV:        "109",
				Account: &models.Account{
					ID:            "account-006",
					AccountNumber: "4444444444",
					Amount:        444.44,
				},
				Customer: &models.Customer{
					ID:        "customer-006",
					FirstName: "William",
					LastName:  "White",
					Email:     "william.white@example.com",
				},
			},
			{
				ID:         "76543",
				CardNumber: "4567-8901-2345-6789",
				HolderName: "Ava Martin",
				ExpDate:    "05/2025",
				CVV:        "765",
				Account: &models.Account{
					ID:            "account-007",
					AccountNumber: "5555555555",
					Amount:        555.55,
				},
				Customer: &models.Customer{
					ID:        "customer-007",
					FirstName: "Ava",
					LastName:  "Martin",
					Email:     "ava.martin@example.com",
				},
			},
			{
				ID:         "54321",
				CardNumber: "6789-0123-4567-8901",
				HolderName: "Isabella Davis",
				ExpDate:    "07/2026",
				CVV:        "432",
				Account: &models.Account{
					ID:            "account-008",
					AccountNumber: "6666666666",
					Amount:        666.66,
				},
				Customer: &models.Customer{
					ID:        "customer-008",
					FirstName: "Isabella",
					LastName:  "Davis",
					Email:     "isabella.davis@example.com",
				},
			},
			{
				ID:         "98765",
				CardNumber: "9012-3456-7890-1234",
				HolderName: "Julian Hall",
				ExpDate:    "04/2025",
				CVV:        "987",
				Account: &models.Account{
					ID:            "account-009",
					AccountNumber: "7777777777",
					Amount:        777.77,
				},
				Customer: &models.Customer{
					ID:        "customer-009",
					FirstName: "Julian",
					LastName:  "Hall",
					Email:     "julian.hall@example.com",
				},
			},
			{
				ID:         "11111",
				CardNumber: "2345-6789-0123-4567",
				HolderName: "Gabriel Brooks",
				ExpDate:    "02/2026",
				CVV:        "654",
				Account: &models.Account{
					ID:            "account-010",
					AccountNumber: "8888888888",
					Amount:        888.88,
				},
				Customer: &models.Customer{
					ID:        "customer-010",
					FirstName: "Gabriel",
					LastName:  "Brooks",
					Email:     "gabriel.brooks@example.com",
				},
			},
		}

		err := h.repo.FillupData(cards)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed filling up data"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data added successfully"))
	}
}

func (h *HttpHandler) GetTransactionHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			h.log.Error("Token is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token is required"))
			return
		}

		statusCode, body, err := h.bankAuthenticator.IsValidToken(token)
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

		accountNumber := r.Header.Get("account_number")
		if accountNumber == "" {
			h.log.Error("Account Number is required")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Account Number is required"))
			return
		}

		transactions, err := h.repo.GetTransactionHistory(accountNumber)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get transaction history"))
			return
		}

		transacionBytes, err := json.Marshal(transactions)
		if err != nil {
			log.WithError(err).Error("Failed to marshal transactions")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal transactions"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(transacionBytes)
	}
}
