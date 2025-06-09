package database

import (
	"context"
	"sort"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) SaveActivityReport(input *model.NewActivityReport) *model.ActivityReport {
	collection := db.client.Database("go_trading_db").Collection("ActivityReports")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error save func:")
	}
	return &model.ActivityReport{
		ID:             res.InsertedID.(primitive.ObjectID).Hex(),
		Timestamp:      input.Timestamp,
		Qty:            input.Qty,
		AvgGain:        input.AvgGain,
		TopAGain:       input.TopAGain,
		TopBGain:       input.TopBGain,
		TopCGain:       input.TopCGain,
		FearGreedIndex: input.FearGreedIndex,
	}
}

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

func (db *DB) FindActivityReportByID(ID string) *model.ActivityReport {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Error().Err(err).Msg("Error find by func:")
	}
	collection := db.client.Database("go_trading_db").Collection("ActivityReports")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	ActivityReport := model.ActivityReport{}
	res.Decode(&ActivityReport)
	return &ActivityReport
}

func (db *DB) AllActivityReports() []*model.ActivityReport {
	collection := db.client.Database("go_trading_db").Collection("ActivityReports")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Error querying database:")
	}

	var ActivityReports []*model.ActivityReport
	for cur.Next(ctx) {
		var ActivityReport model.ActivityReport
		err := cur.Decode(&ActivityReport)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding document:")
		}
		ActivityReports = append(ActivityReports, &ActivityReport)
	}

	return ActivityReports
}

func (db *DB) SaveTradeOutcomeReport(input *model.NewTradeOutcomeReport) *model.TradeOutcomeReport {
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error save func:")
	}
	return &model.TradeOutcomeReport{
		ID:               res.InsertedID.(primitive.ObjectID).Hex(),
		Timestamp:        input.Timestamp,
		BotName:          input.BotName,
		PercentageChange: input.PercentageChange,
		Balance:          input.Balance,
		Symbol:           input.Symbol,
		Outcome:          input.Outcome,
		Fee:              input.Fee,
		ElapsedTime:      input.ElapsedTime,
		Volume:           input.Volume,
		FearGreedIndex:   input.FearGreedIndex,
		MarketStatus:     input.MarketStatus,
	}
}

func (db *DB) FindTradeOutcomeReportByID(ID string) *model.TradeOutcomeReport {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Error().Err(err).Msg("Error find by func:")
	}
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	TradeOutcomeReport := model.TradeOutcomeReport{}
	res.Decode(&TradeOutcomeReport)
	return &TradeOutcomeReport
}

func (db *DB) AllTradeOutcomeReports() []*model.TradeOutcomeReport {
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Error All func:")
	}
	var TradeOutcomeReports []*model.TradeOutcomeReport
	for cur.Next(ctx) {
		var TradeOutcomeReport *model.TradeOutcomeReport
		err := cur.Decode(&TradeOutcomeReport)
		if err != nil {
			log.Error().Err(err).Msg("Error Decode func:")
		}
		TradeOutcomeReports = append(TradeOutcomeReports, TradeOutcomeReport)
	}
	return TradeOutcomeReports
}

func (db *DB) FindSymbolStatsBySymbol(ctx context.Context, symbol string) *model.SymbolStats {
	collection := db.client.Database("go_trading_db").Collection("SymbolStats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"symbol": symbol})
	symbolStats := model.SymbolStats{}
	res.Decode(&symbolStats)
	return &symbolStats
}

func (db *DB) AllSymbolStats() []*model.SymbolStats {
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

// TradeOutcomeReportsByBot retrieves trade outcome reports based on the BotName.
func (db *DB) TradeOutcomeReportsByBotName(ctx context.Context, botName string) ([]*model.TradeOutcomeReport, error) {
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")

	filter := bson.D{{"botname", botName}}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error TradeOutcomeReportsByBot func:")
		return nil, err
	}

	var tradeOutcomeReports []*model.TradeOutcomeReport
	for cur.Next(ctx) {
		var tradeOutcomeReport *model.TradeOutcomeReport
		err := cur.Decode(&tradeOutcomeReport)
		if err != nil {
			log.Error().Err(err).Msg("Error Decode func:")
		}
		tradeOutcomeReports = append(tradeOutcomeReports, tradeOutcomeReport)
	}
	return tradeOutcomeReports, nil
}

// TradeOutcomeReportsByBotNameAndMarketStatus retrieves trade outcome reports based on the BotName and MarketStatus with a limit.
func (db *DB) TradeOutcomeReportsByBotNameAndMarketStatus(ctx context.Context, botName string, marketStatus string, limit int) ([]*model.TradeOutcomeReport, error) {
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")

	filter := bson.D{
		{"botname", botName},
		{"marketstatus", marketStatus},
	}

	findOptions := options.Find()

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error TradeOutcomeReportsByBotName func:")
		return nil, err
	}

	var tradeOutcomeReports []*model.TradeOutcomeReport
	for cur.Next(ctx) {
		var tradeOutcomeReport *model.TradeOutcomeReport
		err := cur.Decode(&tradeOutcomeReport)
		if err != nil {
			log.Error().Err(err).Msg("Error Decode func:")
		}
		tradeOutcomeReports = append(tradeOutcomeReports, tradeOutcomeReport)
	}

	// Sort the data by timestamp in descending order
	sort.Slice(tradeOutcomeReports, func(i, j int) bool {
		return tradeOutcomeReports[i].Timestamp > tradeOutcomeReports[j].Timestamp
	})

	// Apply the limit after sorting
	if limit > 0 && limit < len(tradeOutcomeReports) {
		tradeOutcomeReports = tradeOutcomeReports[:limit]
	}

	return tradeOutcomeReports, nil
}

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

// DeleteStrategy deletes a strategy from the database.
func (db *DB) DeleteTradeOutcomeReport(ctx context.Context, timestamp int) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("TradeOutcomeReports")

	// Define a filter to match documents with the specified timestamp
	filter := bson.D{{"timestamp", timestamp}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting trade outcome from the database:")
		return false, err
	}

	return result.DeletedCount > 0, nil
}
