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

func Connect() (*DB, error) {
	uri := "mongodb://fudgebot:cookiebot@database:27017/go_trading_db"
	log.Info().Str("mongodb_uri", uri).Msg("Connecting to MongoDB")

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Auth = &options.Credential{
		Username:   "fudgebot",
		Password:   "cookiebot",
		AuthSource: "go_trading_db",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to MongoDB")
		return nil, err
	}

	return &DB{client: client}, nil
}

func (db *DB) Close() {
	if db.client != nil {
		db.client.Disconnect(context.Background())
	}
}
