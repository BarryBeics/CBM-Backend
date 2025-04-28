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
