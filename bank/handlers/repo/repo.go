package repo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pavr1/online_payment_platform/bank/config"
	"github.com/pavr1/online_payment_platform/bank/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRepoHandler interface {
	Transfer(fromCard *models.Card, targetAccountNumber string, amount float64, description string) (int, error)
	FillupData(cards []*models.Card) error
}

type RepoHandler struct {
	log    *log.Logger
	Config *config.Config
	client *mongo.Client
}

func NewRepoHandler(log *log.Logger, config *config.Config) (IRepoHandler, error) {
	client, err := connectToMongoDB(config)
	if err != nil {
		log.WithField("error", err).Error("Failed to connect to MongoDB")

		return nil, err
	}

	return &RepoHandler{
		log:    log,
		Config: config,
		client: client,
	}, nil
}

func connectToMongoDB(config *config.Config) (*mongo.Client, error) {
	uri := config.MongoDB.Uri

	log.WithField("uri", uri).Info("Connecting to MongoDB...")

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB")

		return nil, err
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func (r *RepoHandler) loadCardInfo(fieldName, valueName string) (*models.Card, error) {
	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	collection := db.Collection(r.Config.MongoDB.Card_Collection)

	// Find the document by ID
	filter := bson.M{fieldName: valueName}
	var card models.Card
	err := collection.FindOne(context.Background(), filter).Decode(&card)
	if err != nil {
		log.WithFields(log.Fields{"fieldName": fieldName, "valueName": valueName}).WithError(err).Error("Failed to find document in db")

		return nil, err
	}

	return &card, nil
}

func (r *RepoHandler) logTransaction(session mongo.Session, transaction *models.Transaction) error {
	// Insert the person into the "people" collection
	collection := session.Client().Database(r.Config.MongoDB.Database).Collection(r.Config.MongoDB.Transaction_Collection)

	doc := bson.D{}

	// Add fields to the document
	doc = append(doc, bson.E{Key: "id", Value: transaction.ID})
	doc = append(doc, bson.E{Key: "date", Value: transaction.Date})
	doc = append(doc, bson.E{Key: "amount", Value: transaction.Amount})
	doc = append(doc, bson.E{Key: "from_card", Value: transaction.FromCard})
	doc = append(doc, bson.E{Key: "to_account", Value: transaction.ToAccount})
	doc = append(doc, bson.E{Key: "details", Value: transaction.Detail})

	// Convert the document to BSON
	bson, err := bson.Marshal(doc)
	if err != nil {
		log.WithError(err).Error("Failed to marshal person to BSON")
		return err
	}

	_, err = collection.InsertOne(context.Background(), bson)
	if err != nil {
		log.WithError(err).Error("Failed to insert person into MongoDB")

		return err
	}

	log.WithField("id", transaction.ID).Info("Transaction inserted successfully")

	return nil
}

func (r *RepoHandler) startTransaction(fromCard, toCard *models.Card, amount float64, description string) error {
	log.WithFields(log.Fields{
		"fromCard": fromCard,
		"toCard":   toCard,
		"amount":   amount,
	}).Info("Starting Transaction...")

	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	cardCollection := db.Collection(r.Config.MongoDB.Card_Collection)

	session, err := r.client.StartSession()
	if err != nil {
		log.WithError(err).Error("Failed to start transactional session")

		return err
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		log.WithError(err).Error("Failed to start transaction")

		return err
	}

	log.WithField("fromCard", fromCard).Info("Updaing from card...")
	// Update the document by ID
	filter := bson.M{"id": fromCard.ID}
	update := bson.M{"$set": fromCard}
	_, err = cardCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.WithError(err).Error("Failed to update from card in db")

		return err
	}

	log.WithField("toCard", toCard).Info("Updaing to card...")
	// Update the document by ID
	filter = bson.M{"id": toCard.ID}
	update = bson.M{"$set": toCard}
	_, err = cardCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.WithError(err).Error("Failed to update to card in db")

		return err
	}

	transaction := models.Transaction{
		ID:        uuid.New().String(),
		Date:      time.Now(),
		Amount:    amount,
		FromCard:  fromCard.CardNumber,
		ToAccount: toCard.Account.AccountNumber,
		Detail:    description,
	}

	r.logTransaction(session, &transaction)

	log.WithField("id", transaction.ID).Info("Transaction logs inserted successfully")

	// Commit the transaction
	if err := session.CommitTransaction(context.Background()); err != nil {
		log.Error("Failed to commit transaction")
	}

	log.WithFields(log.Fields{
		"fromCard": fromCard,
		"toCard":   toCard,
		"amount":   amount,
	}).Info("Transaction committed")

	return nil
}

func (r *RepoHandler) Transfer(fromCard *models.Card, targetAccountNumber string, amount float64, description string) (int, error) {
	dbFromCard, err := r.loadCardInfo("cardnumber", fromCard.CardNumber)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if fromCard.HolderName != dbFromCard.HolderName {
		log.WithField("HolderName", fromCard.HolderName).Error("Invalid Holder Name")
		return http.StatusBadRequest, fmt.Errorf("invalid request, please check your card information")
	}
	if fromCard.ExpDate != dbFromCard.ExpDate {
		log.WithField("ExpDate", fromCard.ExpDate).Error("Invalid Expiration Date")
		return http.StatusBadRequest, fmt.Errorf("invalid request, please check your card information")
	}
	if fromCard.CVV != dbFromCard.CVV {
		log.WithField("CVV", fromCard.CVV).Error("Invalid CVV")
		return http.StatusBadRequest, fmt.Errorf("invalid request, please check your card information")
	}

	fromCardCurrentAmount := dbFromCard.GetAmount()
	if fromCardCurrentAmount < amount {
		log.WithField("amount", amount).Error("Insufficient balance")
		return http.StatusUnauthorized, fmt.Errorf("invalid request, Insuficient balance")
	}

	dbFromCard.SetAmount(fromCardCurrentAmount - amount)

	dbToCard, err := r.loadCardInfo("account.accountnumber", targetAccountNumber)
	if err != nil {
		return http.StatusBadRequest, err
	}

	dbToCard.SetAmount(dbToCard.GetAmount() + amount)

	err = r.startTransaction(dbFromCard, dbToCard, amount, description)

	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}
	return status, err
}

func (r *RepoHandler) FillupData(cards []*models.Card) error {
	// Insert the person into the "people" collection
	collection := r.client.Database(r.Config.MongoDB.Database).Collection(r.Config.MongoDB.Card_Collection)

	for _, card := range cards {
		_, err := collection.InsertOne(context.Background(), card)
		if err != nil {
			log.WithError(err).Error("Failed to insert card information")

			return err
		}
	}

	return nil
}
