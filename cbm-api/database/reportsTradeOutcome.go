package database

import (
	"context"
	"sort"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTradeOutcomeReport saves a new trade outcome report to the database.
func (db *DB) CreateTradeOutcomeReport(input *model.NewTradeOutcomeReport) *model.TradeOutcomeReport {
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

// ReadAllTradeOutcomeReports retrieves all trade outcome reports from the database.
func (db *DB) ReadAllTradeOutcomes() []*model.TradeOutcomeReport {
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

// ReadTradeOutcomeReportByID retrieves a trade outcome report by its ID from the database.
func (db *DB) ReadTradeOutcomeReportByID(ID string) *model.TradeOutcomeReport {
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

// TradeOutcomeReportsByBot retrieves trade outcome reports based on the BotName.
func (db *DB) ReadTradeOutcomesPerBotName(ctx context.Context, botName string) ([]*model.TradeOutcomeReport, error) {
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
func (db *DB) ReadTradeOutcomeInFocus(ctx context.Context, botName string, marketStatus string, limit int) ([]*model.TradeOutcomeReport, error) {
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

// DeleteStrategy deletes a strategy from the database.
func (db *DB) DeleteOutcomeReports(ctx context.Context, timestamp int) (bool, error) {
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
