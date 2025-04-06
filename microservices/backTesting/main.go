package main

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cryptobotmanager.com/cbm-backend/microservices/priceData/functions"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
	// Add any other necessary imports
)

func main() {
	// Setup GraphQL backend client
	backend := os.Getenv("TRADING_BOT_URL")
	if backend == "" {
		backend = "http://resolvers:8080/query"
	}

	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()

	// Directory with your JSON price files
	dataDir := "../../binancePrices"

	// Read and sort files chronologically
	var files []string
	filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err).Msg("Error walking through files")
			return err
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), "binance_prices_") && strings.HasSuffix(d.Name(), ".json") {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files) // Ensure they’re processed in date order

	// Loop through each file and "replay" the data
	for _, file := range files {
		log.Info().Str("File", file).Msg("Processing file")

		marketData, err := functions.LoadPriceSnapshotsFromFile(file)
		if err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to load JSON data")
			continue
		}

		for _, snapshot := range marketData {
			// Convert []Price (in snapshot.Pairs) to []PriceData
			var market []functions.PriceData
			for _, p := range snapshot.Pairs {
				market = append(market, functions.PriceData{
					Symbol: p.Symbol,
					Price:  p.Price,
				})
			}

			err := functions.SavePriceData(ctx, client, market, int(snapshot.Timestamp))
			if err != nil {
				log.Error().Err(err).Int64("timestamp", snapshot.Timestamp).Msg("Failed to save snapshot")
			} else {
				log.Info().Int64("timestamp", snapshot.Timestamp).Msg("Replayed snapshot")
			}
		}

	}
}
