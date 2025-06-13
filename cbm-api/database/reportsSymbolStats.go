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

// ==========================
// === Ticker Stats ===
// ==========================

// CreateHistoricTickerStats saves historic ticker stats to the database.
func (db *DB) CreateHistoricTickerStats(input model.NewHistoricTickerStatsInput) ([]*model.HistoricTickerStats, error) {
	log.Info().
		Int("timestamp", input.Timestamp).
		Int("num_stats", len(input.Stats)).
		Msg("Preparing to insert historic ticker stats")

	collection := db.client.Database("go_trading_db").Collection("HistoricTickerStats")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var historicTickerStats []*model.HistoricTickerStats
	for _, statInput := range input.Stats {
		historicTickerStats = append(historicTickerStats, &model.HistoricTickerStats{
			Timestamp: input.Timestamp,
			Stats: []*model.TickerStats{{
				Symbol:            statInput.Symbol,
				PriceChange:       statInput.PriceChange,
				PriceChangePct:    statInput.PriceChangePct,
				QuoteVolume:       statInput.QuoteVolume,
				Volume:            statInput.Volume,
				TradeCount:        statInput.TradeCount,
				HighPrice:         statInput.HighPrice,
				LowPrice:          statInput.LowPrice,
				LastPrice:         statInput.LastPrice,
				LiquidityEstimate: statInput.LiquidityEstimate,
			}},
			CreatedAt: time.Now().UTC(),
		})
	}

	// Convert to []interface{} for MongoDB
	docs := make([]interface{}, len(historicTickerStats))
	for i, doc := range historicTickerStats {
		docs[i] = doc
	}

	log.Info().Msg("Calling InsertMany on HistoricTickerStats")

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Error().Err(err).Msg("Error saving historic ticker stats:")
		return nil, err
	}

	return historicTickerStats, nil
}

// ReadTickerStatsBySymbol retrieves ticker stats by symbol from the database.
func (db *DB) ReadTickerStatsBySymbol(symbol string, limit int) ([]*model.TickerStats, error) {
	log.Info().Msgf("Querying ticker stats at Symbol: %s", symbol)
	collection := db.client.Database("go_trading_db").Collection("HistoricTickerStats")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.M{"stats.symbol": symbol}
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching ticker stats by symbol")
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []model.HistoricTickerStats
	if err := cursor.All(ctx, &docs); err != nil {
		log.Error().Err(err).Msg("Failed to decode raw bson")
		return nil, err
	}

	var tickerStats []*model.TickerStats
	for _, doc := range docs {
		for _, stat := range doc.Stats {
			if stat.Symbol == symbol {
				tickerStats = append(tickerStats, stat)
			}
		}
	}

	log.Info().Int("count", len(tickerStats)).Msg("Fetched ticker stats by symbol")
	return tickerStats, nil
}

// ReadHistoricTickerStatsAtTimestamp retrieves historic ticker stats at a specific timestamp.
func (db *DB) ReadHistoricTickerStatsAtTimestamp(timestamp int) ([]model.HistoricTickerStats, error) {
	log.Info().Msgf("Querying historic ticker stats at Timestamp: %d", timestamp)
	collection := db.client.Database("go_trading_db").Collection("HistoricTickerStats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Filter by timestamp
	filter := bson.M{"timestamp": timestamp}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching historic ticker stats at timestamp")
		return nil, err
	}
	defer cursor.Close(ctx)

	var historicTickerStats []model.HistoricTickerStats

	// Iterate over the results
	for cursor.Next(ctx) {
		var result model.HistoricTickerStats
		if err := cursor.Decode(&result); err != nil {
			log.Error().Err(err).Msg("Error decoding historic ticker stats at timestamp")
			return nil, err
		}

		// Append the result to the list
		historicTickerStats = append(historicTickerStats, result)
	}

	return historicTickerStats, nil
}

// DeleteHistoricTickerStatsByTimestamp deletes historic ticker stats by timestamp.
func (db *DB) DeleteHistoricTickerStatsByTimestamp(ctx context.Context, timestamp int) error {
	collection := db.client.Database("go_trading_db").Collection("HistoricTickerStats")

	// Define a filter to match documents with the specified symbol
	filter := bson.D{{"timestamp", timestamp}}

	// Perform the delete operation
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting historic ticker stats with timestamp %s", timestamp)
		return err
	}

	log.Info().Msgf("Deleted %d historic ticker stats with timestamp %s", result.DeletedCount, timestamp)

	return nil
}

// ==========================
// === Symbol Stats ===
// ==========================

// UpsertSymbolStats updates or inserts symbol statistics in the database.
func (db *DB) UpsertSymbolStats(input *model.UpsertSymbolStatsInput) *model.SymbolStats {
	collection := db.client.Database("go_trading_db").Collection("SymbolStats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"symbol": input.Symbol}

	// Fetch current document
	var existing model.SymbolStats
	err := collection.FindOne(ctx, filter).Decode(&existing)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error().Err(err).Msg("Failed to fetch existing symbol stats")
		return nil
	}

	// Merge PositionCounts averages
	mergedCounts := input.PositionCounts
	if len(existing.PositionCounts) == len(input.PositionCounts) {
		for i := range input.PositionCounts {
			old := existing.PositionCounts[i]
			new := input.PositionCounts[i]
			totalCount := old.Count + new.Count
			if totalCount > 0 {
				mergedCounts[i].Avg = ((old.Avg * float64(old.Count)) + (new.Avg * float64(new.Count))) / float64(totalCount)
				mergedCounts[i].Count = totalCount
			}
		}
	}

	// Merge LiquidityEstimate
	var mergedLiquidity *model.Mean
	if input.LiquidityEstimate != nil {
		mergedLiquidity = &model.Mean{
			Avg:   input.LiquidityEstimate.Avg,
			Count: input.LiquidityEstimate.Count,
		}
	}
	if existing.LiquidityEstimate != nil && mergedLiquidity != nil {
		old := existing.LiquidityEstimate
		new := mergedLiquidity
		totalCount := old.Count + new.Count
		if totalCount > 0 {
			mergedLiquidity = &model.Mean{
				Avg:   ((old.Avg * float64(old.Count)) + (new.Avg * float64(new.Count))) / float64(totalCount),
				Count: totalCount,
			}
		}
	}

	update := bson.M{
		"$set": bson.M{
			"symbol":               input.Symbol,
			"positionCounts":       mergedCounts,
			"avgLiquidityEstimate": mergedLiquidity,
			"maxLiquidityEstimate": input.MaxLiquidityEstimate,
			"minLiquidityEstimate": input.MinLiquidityEstimate,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Error().Err(err).Msg("Error in upsert")
		return nil
	}

	// Convert []*model.MeanInput to []*model.Mean
	var positionCounts []*model.Mean
	for _, m := range mergedCounts {
		if m != nil {
			positionCounts = append(positionCounts, &model.Mean{
				Avg:   m.Avg,
				Count: m.Count,
			})
		} else {
			positionCounts = append(positionCounts, nil)
		}
	}

	return &model.SymbolStats{
		Symbol:               input.Symbol,
		PositionCounts:       positionCounts,
		LiquidityEstimate:    mergedLiquidity,
		MaxLiquidityEstimate: input.MaxLiquidityEstimate,
		MinLiquidityEstimate: input.MinLiquidityEstimate,
	}
}

// ReadAllSymbolStats retrieves all symbol statistics from the database.
func (db *DB) ReadAllSymbolStats() []*model.SymbolStats {
	collection := db.client.Database("go_trading_db").Collection("SymbolStats")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Error All func:")
	}

	var symbolStats []*model.SymbolStats
	for cur.Next(ctx) {
		var stat model.SymbolStats
		err := cur.Decode(&stat)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding document:")
		}
		symbolStats = append(symbolStats, &stat)
	}

	return symbolStats
}

// ReadSingleSymbolStatsBySymbol retrieves symbol statistics by symbol from the database.
func (db *DB) ReadSingleSymbolStatsBySymbol(symbol string) *model.SymbolStats {
	collection := db.client.Database("go_trading_db").Collection("SymbolStats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"symbol": symbol})
	symbolStats := model.SymbolStats{}
	res.Decode(&symbolStats)
	return &symbolStats
}

// DeleteSymbolStats deletes symbol statistics by symbol from the database.
func (db *DB) DeleteSymbolStats(ctx context.Context, symbol string) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("SymbolStats")

	// Define a filter to match documents with the specified symbol
	filter := bson.D{{"symbol", symbol}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting symbol stats from the database:")
		return false, err
	}

	return result.DeletedCount > 0, nil
}
