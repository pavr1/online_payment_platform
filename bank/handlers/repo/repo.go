package repo

import (
	"context"
	"time"

	"github.com/pavr1/online_payment_platform/bank/config"
	"github.com/pavr1/online_payment_platform/bank/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRepoHandler interface {
	VerifyCard(cardModel *models.Card) (bool, error)
	LogTransaction(transaction *models.Transaction) error
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

func (r *RepoHandler) VerifyCard(cardModel *models.Card) (bool, error) {
	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	collection := db.Collection(r.Config.MongoDB.Collection)

	// Find the document by ID
	filter := bson.M{"card_number": cardModel.CardNumber}
	var card models.Card
	err := collection.FindOne(context.Background(), filter).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.WithField("card_number", cardModel.CardNumber).Error("Document not found in MongoDB")

			return false, nil
		}

		log.WithError(err).Error("Failed to find document in MongoDB")

		return false, err
	}

	if card.HolderName != cardModel.HolderName {
		log.WithField("holder_name", cardModel.HolderName).Error("Holder name does not match")

		return false, nil
	} else if card.ExpDate != cardModel.ExpDate {
		log.WithField("exp_date", cardModel.ExpDate).Error("Expiration date does not match")

		return false, nil
	} else if card.CVV != cardModel.CVV {
		log.WithField("cvv", cardModel.CVV).Error("CVV does not match")

		return false, nil
	}

	return true, nil
}

func (r *RepoHandler) LogTransaction(transaction *models.Transaction) error {
	// Insert the person into the "people" collection
	collection := r.client.Database(r.Config.MongoDB.Database).Collection(r.Config.MongoDB.Collection)

	doc := bson.D{}

	// Add fields to the document
	doc = append(doc, bson.E{Key: "id", Value: transaction.ID})
	doc = append(doc, bson.E{Key: "date", Value: transaction.Date})
	doc = append(doc, bson.E{Key: "amount", Value: transaction.Amount})
	doc = append(doc, bson.E{Key: "from_account", Value: transaction.FromAccount})
	doc = append(doc, bson.E{Key: "to_account", Value: transaction.ToAccount})
	doc = append(doc, bson.E{Key: "details", Value: transaction.Detail})

	// Convert the document to BSON
	personBSON, err := bson.Marshal(doc)
	if err != nil {
		log.WithError(err).Error("Failed to marshal person to BSON")
		return err
	}

	_, err = collection.InsertOne(context.Background(), personBSON)
	if err != nil {
		log.WithError(err).Error("Failed to insert person into MongoDB")

		return err
	}

	log.WithField("id", transaction.ID).Info("Transaction inserted successfully")

	return nil
}
