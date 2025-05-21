package functions

import (
	"context"
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func CSVPrices(backend string) error {
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()

	// Directory with your JSON price files
	dataDir := "binancePrices"

	// Read and sort files chronologically
	var files []string
	filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err).Msg("Error walking through files")
			return err
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), "binance_prices_") && strings.HasSuffix(d.Name(), ".json") {
			log.Debug().Str("discovered_file", path).Msg("Found matching price file")
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files) // Ensure they’re processed in date order

	// Loop through each file and "replay" the data
	for _, file := range files {
		log.Debug().Str("File", file).Msg("Processing file")

		marketData, err := LoadPriceSnapshotsFromFile(file)
		if err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to load JSON data")
			continue
		}

		for _, snapshot := range marketData {
			var market []model.Pair
			for _, p := range snapshot.Pairs {
				market = append(market, model.Pair{
					Symbol: p.Symbol,
					Price:  p.Price,
				})
			}

			if len(market) == 0 {
				log.Warn().
					Int("timestamp", snapshot.Timestamp).
					Msg("Skipping snapshot due to empty market data")
				continue
			}

			if err := SavePriceData(ctx, client, market, int(snapshot.Timestamp)); err != nil {
				log.Error().Err(err).Int("timestamp", snapshot.Timestamp).Msg("Save PriceData")
			}

			// Start trading
			err := letTrade(ctx, client, market, int(snapshot.Timestamp))
			if err != nil {
				log.Error().Err(err).Int("timestamp", snapshot.Timestamp).Msg("lets trade")
			}
		}
	}
	return nil
}

func BinancePrices(backend string) error {
	// make live API calls to Binance
	// Get the nearest whole 5 minutes & print the current time
	now := time.Now().Unix()
	roundedEpoch := shared.RoundTimeToFiveMinuteInterval(now)
	log.Info().Int64("Executing task at:", now).Int("Rounded time", roundedEpoch).Msg("Time")

	// Create Client & Context
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()
	var market []model.Pair
	var err error

	market, err = FetchPricesFromBinanceAPI()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

	err = SavePriceDataAsJSON(market, int64(roundedEpoch))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save price data to JSON!")
	}

	if err := SavePriceData(ctx, client, market, roundedEpoch); err != nil {
		log.Error().Err(err).Int("timestamp", roundedEpoch).Msg("Save PriceData")
	}

	return letTrade(ctx, client, market, roundedEpoch)
}
