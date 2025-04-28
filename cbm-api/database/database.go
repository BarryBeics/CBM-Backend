package database

import (
	"context"

	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

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
		log.Error().Err(err).Msg("Error client options func:")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error ctx func:")
	}

	return &DB{
		client: client,
	}
}

func (db *DB) Close() {
	if db.client != nil {
		db.client.Disconnect(context.Background())
	}
}
