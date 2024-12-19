package database

import (
	"context"

	"time"

	"cryptobotmanager.com/cbm-backend/Resolvers/graph/model"
	"github.com/rs/zerolog/log"
)

func (db *DB) SaveHistoricPrices(input *model.NewHistoricPriceInput) ([]*model.HistoricPrices, error) {
	log.Info().Msgf("Inserting prices into DB: %+v", input)
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
