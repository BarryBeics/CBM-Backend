package database

import (
	"context"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) UpsertFearAndGreedIndex(ctx context.Context, input model.UpsertFearAndGreedIndexInput) (*model.FearAndGreedIndex, error) {
	collection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")

	filter := bson.M{"timestamp": input.Timestamp}
	update := bson.M{
		"$set": bson.M{
			"timestamp":           input.Timestamp,
			"value":               input.Value,
			"valueClassification": input.ValueClassification,
			"createdAt":           time.Now().UTC(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upsert fear and greed index")
		return nil, err
	}

	// Return the upserted document (optional)
	return &model.FearAndGreedIndex{
		Timestamp:           input.Timestamp,
		Value:               input.Value,
		ValueClassification: input.ValueClassification,
		CreatedAt:           time.Now().UTC(),
	}, nil
}

func (db *DB) DeleteFearAndGreedIndex(ctx context.Context, timestamp int) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")

	filter := bson.M{"timestamp": timestamp}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete fear and greed index")
		return false, err
	}

	if result.DeletedCount == 0 {
		log.Warn().Int("timestamp", timestamp).Msg("No fear and greed index found to delete")
		return false, nil
	}

	return true, nil
}

func (db *DB) ReadFearAndGreedIndex(ctx context.Context, limit *int) ([]*model.FearAndGreedIndex, error) {
	collection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")

	var results []*model.FearAndGreedIndex
	findOptions := options.Find().SetSort(bson.D{{"timestamp", -1}})
	if limit != nil {
		findOptions.SetLimit(int64(*limit))
	}

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find fear and greed index")
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var index model.FearAndGreedIndex
		if err := cursor.Decode(&index); err != nil {
			log.Error().Err(err).Msg("Failed to decode fear and greed index")
			return nil, err
		}
		results = append(results, &index)
	}

	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("Cursor error")
		return nil, err
	}

	return results, nil
}

func (db *DB) ReadFearAndGreedIndexAtTimestamp(ctx context.Context, timestamp int) (*model.FearAndGreedIndex, error) {
	collection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")

	filter := bson.M{"timestamp": timestamp}
	var index model.FearAndGreedIndex

	err := collection.FindOne(ctx, filter).Decode(&index)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Warn().Int("timestamp", timestamp).Msg("No fear and greed index found at this timestamp")
			return nil, nil
		}
		log.Error().Err(err).Msg("Failed to find fear and greed index at timestamp")
		return nil, err
	}

	return &index, nil
}

func (db *DB) ReadFearAndGreedIndexCount(ctx context.Context) (int, error) {
	collection := db.client.Database("go_trading_db").Collection("fear_and_greed_index")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to count fear and greed index documents")
		return 0, err
	}

	return int(count), nil
}
