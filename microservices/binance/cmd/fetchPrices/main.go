package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	backTesting "cryptobotmanager.com/cbm-backend/microservices/backTesting/functions"
	"cryptobotmanager.com/cbm-backend/microservices/binance"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/Khan/genqlient/graphql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: No .env file found or failed to load")
	}

	// Setup GraphQL backend client
	backend := os.Getenv("TRADING_BOT_URL")
	if backend == "" {
		backend = "http://cbm-api:8080/query"
	}
	// Initialize logger
	shared.SetupLogger()

	err = BinancePrices(backend)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

}

// FetchPricesFromBinanceAPI fetches market prices from Binance API
// using API and Secret keys from environment variables.
// It returns a slice of PriceData structs and an error if any.
func FetchPricesFromBinanceAPI() (market []model.Pair, err error) {
	client := binance.NewBinanceClient()

	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("NewListPricesService")
		return nil, err
	}
	fmt.Println("qty", len(prices))

	// iterate over prices and build into slice of structs
	for _, price := range prices {
		market = append(market, model.Pair{Symbol: price.Symbol, Price: price.Price})
		log.Info().Str("Symbol:", price.Price).Str("Price", price.Price)
	}

	return market, nil
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

	err = shared.SavePriceDataAsJSON(market, int64(roundedEpoch))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save price data to JSON!")
	}

	if err := shared.SavePriceData(ctx, client, market, roundedEpoch); err != nil {
		log.Error().Err(err).Int("timestamp", roundedEpoch).Msg("Save PriceData")
	}
	return backTesting.LetsTrade(ctx, client, market, roundedEpoch)
}
