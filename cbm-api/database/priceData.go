package database

import (
	"context"
	"errors"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) SaveHistoricPrices(input *model.NewHistoricPriceInput) ([]*model.HistoricPrices, error) {
	log.Debug().Msgf("Inserting prices into DB: %+v", input)
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a slice to store the inserted HistoricPrices
	var insertedHistoricPrices []*model.HistoricPrices

	// Iterate over pairs and insert each one into the collection
	for _, pairInput := range input.Pairs {
		// Create a new HistoricPrices object for each pair with the provided timestamp
		historicPrices := &model.HistoricPrices{
			Pair:      []*model.Pair{{Symbol: pairInput.Symbol, Price: pairInput.Price}},
			Timestamp: input.Timestamp,
		}

		// Insert the new HistoricPrices object into the collection
		_, err := collection.InsertOne(ctx, historicPrices)
		if err != nil {
			log.Error().Err(err).Msg("Error saving historic price:")
			// Handle the error, perhaps return an error or log it
			return nil, err
		}

		// Append the inserted HistoricPrices to the result slice
		insertedHistoricPrices = append(insertedHistoricPrices, historicPrices)
	}

	// Return the array of inserted HistoricPrices
	return insertedHistoricPrices, nil
}

// HistoricPricesBySymbol fetches historic prices based on the given symbol and limit.
func (db *DB) HistoricPricesBySymbol(symbol string, limit int, ascending bool) ([]model.HistoricPrices, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricPrices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"Pair.Symbol", symbol}}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	sortOrder := -1
	if ascending {
		sortOrder = 1
	}
	findOptions.SetSort(bson.D{{Key: "Timestamp", Value: sortOrder}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error querying historic prices by symbol")
		return nil, err
	}
	defer cursor.Close(ctx)

	var historicPrices []model.HistoricPrices
	if err := cursor.All(ctx, &historicPrices); err != nil {
		log.Error().Err(err).Msg("Error decoding historic prices by symbol")
		return nil, err
	}

	return historicPrices, nil
}

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

// HistoricPricesAtTimestamp fetches historic prices at a specific timestamp.
func (db *DB) HistoricPricesAtTimestamp(timestamp int) ([]model.HistoricPrices, error) {
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

// GetUniqueTimestampCount fetches the count of unique timestamps.
func (db *DB) GetUniqueTimestampCount(ctx context.Context) (int, error) {
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

func (db *DB) SaveHistoricKlineData(input *model.NewHistoricKlineDataInput) ([]*model.HistoricKlineData, error) {
	if input == nil {
		log.Error().Msg("Input is nil")
		return nil, errors.New("input cannot be nil")
	}

	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var insertedKlineData []*model.HistoricKlineData

	// Validate input fields
	if input.Opentime == 0 {
		log.Error().Msg("Opentime is missing")
		return nil, errors.New("opentime is required")
	}
	if len(input.Coins) == 0 {
		log.Error().Msg("Coins list is empty")
		return nil, errors.New("coins list cannot be empty")
	}

	var ohlcs []*model.Ohlc
	for _, coinInput := range input.Coins {
		if coinInput.Symbol == "" {
			log.Error().Msg("Coin symbol is missing")
			return nil, errors.New("coin symbol is required")
		}

		// Convert OHLCInput to OHLC
		ohlc := &model.Ohlc{
			Symbol:      coinInput.Symbol,
			OpenPrice:   coinInput.OpenPrice,
			HighPrice:   coinInput.HighPrice,
			LowPrice:    coinInput.LowPrice,
			ClosePrice:  coinInput.ClosePrice,
			TradeVolume: coinInput.TradeVolume,
		}

		ohlcs = append(ohlcs, ohlc)
	}

	// Create a new HistoricKlineData object with the provided input
	historicKlineData := &model.HistoricKlineData{
		Opentime: input.Opentime,
		Coins:    ohlcs,
	}

	log.Info().Msgf("Saving historic kline data: %+v", historicKlineData)

	// Insert the new HistoricKlineData object into the collection
	_, err := collection.InsertOne(ctx, historicKlineData)
	if err != nil {
		log.Error().Err(err).Msg("Error saving historic kline data")
		return nil, err
	}

	insertedKlineData = append(insertedKlineData, historicKlineData)

	log.Info().Msgf("Successfully saved historic kline data: %+v", insertedKlineData)

	return insertedKlineData, nil
}

func (db *DB) HistoricKlineDataBySymbol(symbol string, limit int) ([]model.HistoricKlineData, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"coins.symbol": symbol} // Assuming your data model has a nested "coins" field

	// Sort the results in descending order based on the opentime field.
	sort := options.Find().SetSort(bson.D{{"opentime", -1}})

	cursor, err := collection.Find(ctx, filter, sort, options.Find().SetLimit(int64(limit)))
	if err != nil {
		log.Error().Err(err).Msg("Error fetching historic kline data by symbol")
		return nil, err
	}
	defer cursor.Close(ctx)

	var klineData []model.HistoricKlineData
	if err := cursor.All(ctx, &klineData); err != nil {
		log.Error().Err(err).Msg("Error decoding historic kline data")
		return nil, err
	}

	return klineData, nil
}

func (db *DB) HistoricKlineDataAtOpentime(opentime int) ([]model.HistoricKlineData, error) {
	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Filter by opentime
	filter := bson.M{"opentime": opentime}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching historic kline data at opentime")
		return nil, err
	}
	defer cursor.Close(ctx)

	var klineData []model.HistoricKlineData

	// Iterate over the results
	for cursor.Next(ctx) {
		var result model.HistoricKlineData
		if err := cursor.Decode(&result); err != nil {
			log.Error().Err(err).Msg("Error decoding historic kline data at opentime")
			return nil, err
		}

		// Append the result to the list
		klineData = append(klineData, result)
	}

	return klineData, nil
}

func (db *DB) DeleteHistoricKlineDataByOpentime(ctx context.Context, opentime int) error {
	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")

	// Define a filter to match documents with the specified opentime
	filter := bson.D{{"opentime", opentime}}

	// Perform the delete operation
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting historic kline data with opentime %d", opentime)
		return err
	}

	log.Info().Msgf("Deleted %d historic kline data with opentime %d", result.DeletedCount, opentime)

	return nil
}
