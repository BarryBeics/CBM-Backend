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

// CreateHistoricKlineData saves historic kline data to the database.
func (db *DB) CreateHistoricKlineData(input *model.NewHistoricKlineDataInput) ([]*model.HistoricKlineData, error) {
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

// ReadHistoricKlineDataBySymbol retrieves historic kline data for a specific symbol.
func (db *DB) ReadHistoricKlineDataBySymbol(symbol string, limit int) ([]model.HistoricKlineData, error) {
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

// MAY NEED TO BE DELETED WILL BECOME CLEAR WHEN WE GET TO KLINE DATA
// // HistoricKlineDataAtOpentime retrieves historic kline data at a specific opentime.
// func (db *DB) HistoricKlineDataAtOpentime(opentime int) ([]model.HistoricKlineData, error) {
// 	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Filter by opentime
// 	filter := bson.M{"opentime": opentime}

// 	cursor, err := collection.Find(ctx, filter)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error fetching historic kline data at opentime")
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var klineData []model.HistoricKlineData

// 	// Iterate over the results
// 	for cursor.Next(ctx) {
// 		var result model.HistoricKlineData
// 		if err := cursor.Decode(&result); err != nil {
// 			log.Error().Err(err).Msg("Error decoding historic kline data at opentime")
// 			return nil, err
// 		}

// 		// Append the result to the list
// 		klineData = append(klineData, result)
// 	}

// 	return klineData, nil
// }

// // DeleteHistoricKlineDataByOpentime deletes historic kline data by opentime.
// func (db *DB) DeleteHistoricKlineDataByOpentime(ctx context.Context, opentime int) error {
// 	collection := db.client.Database("go_trading_db").Collection("HistoricKlineData")

// 	// Define a filter to match documents with the specified opentime
// 	filter := bson.D{{"opentime", opentime}}

// 	// Perform the delete operation
// 	result, err := collection.DeleteMany(ctx, filter)
// 	if err != nil {
// 		log.Error().Err(err).Msgf("Error deleting historic kline data with opentime %d", opentime)
// 		return err
// 	}

// 	log.Info().Msgf("Deleted %d historic kline data with opentime %d", result.DeletedCount, opentime)

// 	return nil
// }
