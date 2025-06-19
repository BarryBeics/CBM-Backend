package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	binance "cryptobotmanager.com/cbm-backend/microservices/externalDataAPIs"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
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

	err = BinanceDailyLiquiditySnapshot(backend)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

}

type QuoteAssetPrices map[string]string

func BinanceDailyLiquiditySnapshot(backend string) error {
	// STEP 1: Get start-of-day timestamp (rounded to 00:00 UTC)
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	log.Info().Int64("timestamp", startOfDay).Msg("Starting daily liquidity snapshot")

	// STEP 2: Create GraphQL client and context
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()

	// STEP 3: Fetch 24h stats from Binance
	market, err := Fetch24hrTickerStats()
	if err != nil || len(market) == 0 {
		log.Error().Err(err).Msg("Failed to get 24h stats from Binance")
		return err
	}

	// STEP 4: Fetch quote asset USD prices (e.g., BNB, BTC, ETH)
	quotePrices, err := GetUSDPricesForQuoteAssets(ctx, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch quote prices")
		return err
	}

	for asset, price := range quotePrices {
		log.Info().Str("quote", asset).Str("usd", price).Msg("Quote asset USD price")
	}

	// STEP 5: Calculate liquidity estimate from 24h volume and trade count
	for i, stat := range market {
		usdVolume := EstimateUSDVolume(stat.Symbol, stat.QuoteVolume, quotePrices)
		if stat.TradeCount <= 0 || usdVolume <= 0 {
			continue // skip
		}

		// LiquidityEstimate = avg USD volume per 5min = (usdVol / trades) / 12
		liquidityEstimate := (usdVolume / float64(stat.TradeCount)) / 12
		str := fmt.Sprintf("%f", liquidityEstimate)
		market[i].LiquidityEstimate = &str
	}

	// Log sample
	for i, s := range market {
		if i >= 10 {
			break
		}
		liq := "nil"
		if s.LiquidityEstimate != nil {
			liq = *s.LiquidityEstimate
		}
		log.Info().Str("symbol", s.Symbol).Str("liq", liq).Msg("Final estimate")
	}

	// STEP 6: Save stats to DB and local JSON
	err = shared.SaveTradeStatsAsJSON(market, startOfDay)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save trade stats as JSON")
	}

	if err := shared.SaveTradeStats(ctx, client, market, int(startOfDay)); err != nil {
		log.Error().Err(err).Msg("Failed to save to DB")
		return err
	}

	return nil
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

	// for _, s := range stats {
	// 	fmt.Printf("Symbol: %s, LastPrice: %s, Volume: %s\n", s.Symbol, s.LastPrice, s.Volume)
	// }

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

func GetUSDPricesForQuoteAssets(ctx context.Context, client graphql.Client) (QuoteAssetPrices, error) {
	requiredSymbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	usdPrices := make(QuoteAssetPrices)

	for _, symbol := range requiredSymbols {
		resp, err := graph.ReadHistoricPrice(ctx, client, symbol, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch price for %s: %w", symbol, err)
		}
		if len(resp.ReadHistoricPrice) == 0 || len(resp.ReadHistoricPrice[0].Pair) == 0 {
			return nil, fmt.Errorf("no price data for symbol: %s", symbol)
		}

		// Extract quote asset (e.g., BTCUSDT → BTC)
		quoteAsset := strings.TrimSuffix(symbol, "USDT")
		usdPrices[quoteAsset] = resp.ReadHistoricPrice[0].Pair[0].Price
	}

	return usdPrices, nil
}

func EstimateUSDVolume(symbol string, quoteVolStr string, quoteAssetPrices QuoteAssetPrices) float64 {
	// Extract quote asset from symbol (e.g., ETHBTC → BTC)
	var quoteAsset string
	for asset := range quoteAssetPrices {
		if strings.HasSuffix(symbol, asset) {
			quoteAsset = asset
			break
		}
	}
	if quoteAsset == "" {
		return 0.0 // unknown quote asset
	}

	quoteVol, err := strconv.ParseFloat(quoteVolStr, 64)
	if err != nil {
		return 0.0 // fallback to 0 if conversion fails
	}

	usdPriceStr := quoteAssetPrices[quoteAsset]
	usdPrice, err := strconv.ParseFloat(usdPriceStr, 64)
	if err != nil {
		return 0.0 // fallback to 0 if conversion fails
	}
	return quoteVol * usdPrice
}
