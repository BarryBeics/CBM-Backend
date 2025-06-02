package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
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

	err = BinanceTradeVolumes(backend)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

}

// FetchPricesFromBinanceAPI fetches market prices from Binance API
// using API and Secret keys from environment variables.
// It returns a slice of PriceData structs and an error if any.
// Fetch24hrTickerStats fetches 24hr stats from Binance API
func Fetch24hrTickerStats() ([]model.TickerStatsInput, error) {
	client := binance.NewBinanceClient() // Make sure this wraps proper API key usage

	stats, err := client.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch price stats")
	}

	for _, s := range stats {
		fmt.Printf("Symbol: %s, LastPrice: %s, Volume: %s\n", s.Symbol, s.LastPrice, s.Volume)
	}

	var results []model.TickerStatsInput
	for _, s := range stats {
		results = append(results, model.TickerStatsInput{
			Symbol:         s.Symbol,
			PriceChange:    s.PriceChange,
			PriceChangePct: s.PriceChangePercent,
			QuoteVolume:    s.QuoteVolume,
			Volume:         s.Volume,
			TradeCount:     int(s.Count),
			HighPrice:      s.HighPrice,
			LowPrice:       s.LowPrice,
			LastPrice:      s.LastPrice,
		})
	}

	return results, nil
}

func BinanceTradeVolumes(backend string) error {
	// make live API calls to Binance
	// Get the nearest whole 5 minutes & print the current time
	now := time.Now().Unix()
	roundedEpoch := shared.RoundTimeToFiveMinuteInterval(now)
	log.Info().Int64("Executing task at:", now).Int("Rounded time", roundedEpoch).Msg("Time")

	// Create Client & Context
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()
	var market []model.TickerStatsInput
	var err error

	market, err = Fetch24hrTickerStats()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

	err = shared.SaveTradeStatsAsJSON(market, int64(roundedEpoch))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save trade stats to JSON!")
	}

	if err := shared.SaveTradeStats(ctx, client, market, roundedEpoch); err != nil {
		log.Error().Err(err).Int("timestamp", roundedEpoch).Msg("Save TradeStats")
	}
	// return backTesting.LetsTrade(ctx, client, market, roundedEpoch)
	return nil
}
