package database

import (
	"context"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// / CreateActivityReport saves a new activity report to the database.
func (db *DB) CreateActivityReport(input *model.NewActivityReport) *model.ActivityReport {
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

// TODO - Activity report not being generated with an ID, need to check why
// ReadActivityReportByID retrieves an activity report by its ID from the database.
func (db *DB) ReadActivityReportByID(ID string) *model.ActivityReport {
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

// ReadAllActivityReports retrieves all activity reports from the database.
func (db *DB) ReadAllActivityReports() []*model.ActivityReport {
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
