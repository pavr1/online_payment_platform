package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pavr1/online_payment_platform/auth/config"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	bankSecretKey            []byte
	paymentPlatformSecretKey []byte
	log                      *log.Logger
	config                   *config.Config
}

func NewHandler(log *log.Logger, secretKey, paymentPlatformSecretKey []byte) *Handler {
	return &Handler{
		bankSecretKey:            secretKey,
		paymentPlatformSecretKey: paymentPlatformSecretKey,
		log:                      log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Info("Handling GET request")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			log.Warn("Missing authorization header")
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err := h.verifyToken(tokenString, h.bankSecretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err.Error())
			log.Warn(err)
			return
		}

		log.Info("Token verified")
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPost {
		log.Info("Handling POST request")
		userName := r.Header.Get("X-User-Name")

		if userName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing X-User-Name"))
			log.Warn("Missing X-User-Name")
			return
		}
		entityName := r.Header.Get("X-Entity-Name")

		if entityName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing X-Entity-Name"))
			log.Warn("Missing X-Entity-Name")
			return
		}

		var secretKey []byte = nil

		if entityName == "PaymentPlatform" {
			secretKey = h.paymentPlatformSecretKey
		} else if entityName == "Bank" {
			secretKey = h.bankSecretKey
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid entity name"))
			log.Warn("Invalid entity name")
			return
		}

		entityKey := r.Header.Get("X-Entity-Key")

		if entityKey == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing X-Entity-Key"))
			log.Warn("Missing X-Entity-Key")
			return
		}

		if entityKey != string(secretKey) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("invalid entity key"))
			log.Warn("Invalid entity key")
			return
		}

		token, err := h.createToken(userName, entityName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to create token"))
			log.Error(err)
			return
		}

		log.Info("Token created")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	} else {
		log.Info("Handling unsupported request")

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createToken(username, entityName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username":   username,
			"entityName": entityName,
			"exp":        time.Now().Add(time.Minute * 5).Unix(),
		})

	tokenString, err := token.SignedString(h.bankSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) verifyToken(tokenString string, secretKey []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
