package database

import (
	"context"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateHistoricPrices saves historic prices to the database.
func (db *DB) CreateHistoricPrices(input *model.NewHistoricPriceInput) ([]*model.HistoricPrices, error) {
	log.Info().
		Int("timestamp", input.Timestamp).
		Int("num_pairs", len(input.Pairs)).
		Msg("Preparing to insert historic prices")

	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	historicPrices := &model.HistoricPrices{
		Pair:      make([]*model.Pair, len(input.Pairs)),
		Timestamp: input.Timestamp,
		CreatedAt: time.Now().UTC(),
	}

	for i, pairInput := range input.Pairs {
		historicPrices.Pair[i] = &model.Pair{
			Symbol:           pairInput.Symbol,
			Price:            pairInput.Price,
			PercentageChange: pairInput.PercentageChange,
		}
	}

	// Log BSON size before insert
	docBytes, err := bson.Marshal(historicPrices)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal document to BSON")
	} else {
		log.Info().
			Int("bson_size_bytes", len(docBytes)).
			Msg("BSON document size before insert")
	}

	log.Info().Msg("Calling InsertOne on HistoricPrices")

	_, err = collection.InsertOne(ctx, historicPrices)
	if err != nil {
		log.Error().Err(err).Msg("Error saving historic price:")
		return nil, err
	}

	insertedHistoricPrices := []*model.HistoricPrices{historicPrices}
	return insertedHistoricPrices, nil
}

// HistoricPricesBySymbol fetches historic prices based on the given symbol and limit.
func (db *DB) ReadHistoricPricesBySymbol(symbol string, limit int) ([]*model.HistoricPrices, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Fetch documents sorted by timestamp descending
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{
		"pair.symbol": symbol, // Only fetch documents containing the symbol
	}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rawResults []*model.HistoricPrices
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, err
	}

	// Filter to return ONLY the matched symbol in the response
	var filteredResults []*model.HistoricPrices
	for _, entry := range rawResults {
		for _, pair := range entry.Pair {
			if pair.Symbol == symbol {
				filteredResults = append(filteredResults, &model.HistoricPrices{
					Timestamp: entry.Timestamp,
					Pair:      []*model.Pair{pair}, // Only include matched pair
				})
				break
			}
		}
	}

	return filteredResults, nil

}

// ReadHistoricPricesAtTimestamp fetches historic prices at a specific timestamp.
func (db *DB) ReadHistoricPricesAtTimestamp(timestamp int) ([]model.HistoricPrices, error) {
	log.Info().Msgf("Querying prices from DB at Timestamp: %d", timestamp)
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Filter by timestamp
	filter := bson.M{"timestamp": timestamp}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching historic prices at timestamp")
		return nil, err
	}
	defer cursor.Close(ctx)

	var historicPrices []model.HistoricPrices

	// Iterate over the results
	for cursor.Next(ctx) {
		var result model.HistoricPrices
		if err := cursor.Decode(&result); err != nil {
			log.Error().Err(err).Msg("Error decoding historic prices at timestamp")
			return nil, err
		}

		// Append the result to the list
		historicPrices = append(historicPrices, result)
	}

	return historicPrices, nil
}

// ReadUniqueTimestampCount fetches the count of unique timestamps.
func (db *DB) ReadUniqueTimestampCount(ctx context.Context) (int, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")

	// Use aggregation to get unique timestamps
	pipeline := bson.A{
		bson.D{{"$group", bson.D{{"_id", "$timestamp"}}}},
		bson.D{{"$group", bson.D{{"_id", nil}, {"count", bson.D{{"$sum", 1}}}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Error().Err(err).Msg("Error counting unique timestamps")
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Error().Err(err).Msg("Error decoding unique timestamps")
		return 0, err
	}

	// Extract count from the result
	count := 0
	if len(result) > 0 {
		// Ensure that "count" field exists in the result
		if countValue, found := result[0]["count"]; found {
			// Use type assertion to handle both int and int32
			switch v := countValue.(type) {
			case int:
				count = v
			case int32:
				count = int(v)
			default:
				log.Error().Msgf("Unexpected type for count: %T", v)
			}
		}
	}

	return count, nil
}

// ReadAvailableSymbols fetches the distinct list of trading symbols.
func (db *DB) ReadAvailableSymbols() ([]string, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Distinct query on the "pair.symbol" field
	symbols, err := collection.Distinct(ctx, "pair.symbol", bson.M{})
	if err != nil {
		log.Error().Err(err).Msg("Error fetching distinct symbols")
		return nil, err
	}

	// Convert to []string
	var result []string
	for _, s := range symbols {
		if str, ok := s.(string); ok {
			result = append(result, str)
		}
	}

	return result, nil
}

// DeleteHistoricPricesByTimestamp deletes historic prices by the specified timestamp.
func (db *DB) DeleteHistoricPricesByTimestamp(ctx context.Context, timestamp int) error {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")

	// Define a filter to match documents with the specified timestamp
	filter := bson.D{{"timestamp", timestamp}}

	// Perform the delete operation
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting historic prices with timestamp %s", timestamp)
		return err
	}

	log.Info().Msgf("Deleted %d historic prices with timestamp %s", result.DeletedCount, timestamp)

	return nil
}

// AllHistoricPrices fetches all historic prices with optional limit and sorting.
// DONT THINK THIS IS USED
func (db *DB) AllHistoricPrices(limit int, ascending bool) ([]model.HistoricPrices, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	sortOrder := -1
	if ascending {
		sortOrder = 1
	}
	findOptions.SetSort(bson.D{{Key: "Timestamp", Value: -sortOrder}})

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error querying historic prices:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var historicPrices []model.HistoricPrices
	if err := cursor.All(ctx, &historicPrices); err != nil {
		log.Error().Err(err).Msg("Error decoding historic prices:")
		return nil, err
	}

	return historicPrices, nil
}
