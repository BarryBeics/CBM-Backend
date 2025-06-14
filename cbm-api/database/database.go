package database

import (
	"context"

	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

// Connect establishes a connection to the MongoDB database and returns a DB instance.
// It also ensures that the necessary indexes are created for the HistoricPrices collection.
func Connect() *DB {
	uri := "mongodb://fudgebot:cookiebot@database:27017/go_trading_db"
	log.Info().Str("mongodb_uri", uri).Msg("Connecting to MongoDB")

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Auth = &options.Credential{
		Username:   "fudgebot",
		Password:   "cookiebot",
		AuthSource: "admin",
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create MongoDB client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	db := &DB{client: client}

	// 🧠 Ensure indexes are present
	if err := db.ensureIndexes(); err != nil {
		log.Error().Err(err).Msg("Failed to create indexes for HistoricPrices")
	}

	return db
}

// ensureIndexes creates necessary indexes for the HistoricPrices collection.
func (db *DB) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Ensure HistoricPrices indexes
	historicPrices := db.client.Database("go_trading_db").Collection("HistoricPrices")
	historicIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "pair.symbol", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().SetName("symbol_timestamp_desc"),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("timestamp_desc"),
		},
	}
	if _, err := historicPrices.Indexes().CreateMany(ctx, historicIndexes); err != nil {
		return err
	}

	// Ensure FearAndGreedIndex index on timestamp (unique)
	fngCollection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")
	fngIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "timestamp", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("timestamp_unique"),
	}
	if _, err := fngCollection.Indexes().CreateOne(ctx, fngIndex); err != nil {
		return err
	}

	return nil
}

// Close disconnects the MongoDB client.
func (db *DB) Close() {
	if db.client != nil {
		db.client.Disconnect(context.Background())
	}
}
