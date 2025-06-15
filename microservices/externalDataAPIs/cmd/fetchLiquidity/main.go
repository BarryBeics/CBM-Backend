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

	err = BinanceTradeVolumes(backend)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

}

type QuoteAssetPrices map[string]string

func BinanceTradeVolumes(backend string) error {

	// STEP 1 : Get the nearest whole 5 minutes & print the current time
	// STEP 1 : Get the nearest whole 5 minutes & print the current time
	now := time.Now().Unix()
	roundedEpoch := shared.RoundTimeToFiveMinuteInterval(now)
	log.Info().Int64("Executing task at:", now).Int("Rounded time", roundedEpoch).Msg("Time")

	// STEP 2 : Create Client & Context
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()
	var market []model.TickerStatsInput
	var err error

	// STEP 3 : Get previous time (60 minutes ago) from roundedEpoch
	returnedPreviousTime, err := shared.GetPreviousTime(roundedEpoch, 60)
	if err != nil {
		log.Error().Msgf("Failed to get previous time!")
		return err
	}
	previousTime := int(returnedPreviousTime)
	log.Info().Int("Previous time", previousTime).Msg("Previous time")

	// STEP 4 : Fetch previous stats from DB
	previousStats, err := getTradeStatsFromDB(ctx, client, int64(previousTime))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get previous trade stats from DB!")
		return err
	}

	for i, s := range previousStats {
		if i >= 10 {
			break
		}
		log.Info().Str("Symbol", s.Symbol).
			Str("LastPrice", s.LastPrice).
			Str("Volume", s.Volume).
			Str("QuoteVolume", s.QuoteVolume).Msg("Previous Stats")
	}

	// STEP 5 : Retrieve latest stats from Binance
	market, err = Fetch24hrTickerStats()
	if err != nil || len(market) == 0 {
		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	}

	// STEP 6 : If previous stats are not found, we can skip to STEP 9
	if previousStats != nil {
		// STEP 7 : Get latest quote prices from DB
		quotePrices, err := GetUSDPricesForQuoteAssets(ctx, client)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get quote prices for assets!")
			return err
		}

		for asset, price := range quotePrices {
			log.Info().Str("QuoteAsset", asset).Str("USDPrice", price).Msg("Quote Asset USD Price")
		}

		// STEP 8 : range market to calculate liquidity estimates

		// Convert previousStats ([]model.TickerStats) to []model.TickerStatsInput
		var previousStatsInput []model.TickerStatsInput
		for _, s := range previousStats {
			previousStatsInput = append(previousStatsInput, model.TickerStatsInput{
				Symbol:            s.Symbol,
				PriceChange:       s.PriceChange,
				PriceChangePct:    s.PriceChangePct,
				QuoteVolume:       s.QuoteVolume,
				Volume:            s.Volume,
				TradeCount:        s.TradeCount,
				HighPrice:         s.HighPrice,
				LowPrice:          s.LowPrice,
				LastPrice:         s.LastPrice,
				LiquidityEstimate: s.LiquidityEstimate,
			})
		}

		for i, stat := range market {
			usdVolume := EstimateUSDVolume(stat.Symbol, stat.QuoteVolume, quotePrices)
			// if usdVolume > 0 {
			// 	log.Info().Str("Symbol", stat.Symbol).
			// 		Float64("USDVolume", usdVolume).
			// 		Msg("Estimated USD Volume")
			// }

			previous := findPreviousStat(stat.Symbol, previousStatsInput)
			if previous == nil {
				//log.Warn().Str("Symbol", stat.Symbol).Msg("No previous stat found")
				continue
			}

			previousUSDVol := EstimateUSDVolume(previous.Symbol, previous.QuoteVolume, quotePrices)
			deltaVol := usdVolume - previousUSDVol
			deltaTrades := stat.TradeCount - previous.TradeCount

			log.Info().Str("Symbol", stat.Symbol).
				Float64("DeltaVolume", deltaVol).
				Int("DeltaTrades", deltaTrades).
				Msg("Delta stats")

			// Add sanity check to ensure both deltas are positive
			if deltaTrades > 0 && deltaVol > 0 {
				liquidityEstimate := (deltaVol / float64(deltaTrades)) / 12
				s := fmt.Sprintf("%f", liquidityEstimate)
				market[i].LiquidityEstimate = &s

				log.Info().
					Str("Symbol", stat.Symbol).
					Float64("LiquidityEstimate", liquidityEstimate).
					Msg("Calculated liquidity estimate")
			} else {
				// Log a warning when estimate is skipped
				log.Warn().
					Str("Symbol", stat.Symbol).
					Float64("DeltaVolume", deltaVol).
					Int("DeltaTrades", deltaTrades).
					Msg("Skipping liquidity estimate due to non-positive deltas")
			}
		}

	}

	// Log the first 10 entries before saving
	for i, s := range market {
		if i >= 10 {
			break
		}
		liq := "nil"
		if s.LiquidityEstimate != nil {
			liq = *s.LiquidityEstimate
		}
		log.Info().
			Str("Symbol", s.Symbol).
			Str("LastPrice", s.LastPrice).
			Str("Volume", s.Volume).
			Str("LiquidityEstimate", liq).
			Msg("About to save TickerStatsInput")
	}

	// STEP 9 : Save the trade stats to DB and JSON including the calculated liquidity estimates
	err = shared.SaveTradeStatsAsJSON(market, int64(roundedEpoch))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save trade stats to JSON!")
	}

	if err := shared.SaveTradeStats(ctx, client, market, roundedEpoch); err != nil {
		log.Error().Err(err).Int("timestamp", roundedEpoch).Msg("Save TradeStats")
	}
	return nil
}

func findPreviousStat(symbol string, stats []model.TickerStatsInput) *model.TickerStatsInput {
	for _, s := range stats {
		if s.Symbol == symbol {
			return &s
		}
	}
	return nil
}

func getTradeStatsFromDB(ctx context.Context, client graphql.Client, roundedEpoch int64) ([]model.TickerStats, error) {
	resp, err := graph.ReadHistoricTickerStatsAtTimestamp(ctx, client, int(roundedEpoch))
	if err != nil {
		log.Error().Err(err).Int64("timestamp", roundedEpoch).Msg("Failed to get trade stats from DB")
		return nil, err
	}

	if resp == nil || len(resp.ReadHistoricTickerStatsAtTimestamp) == 0 {
		log.Warn().Int64("timestamp", roundedEpoch).Msg("No trade stats found for the specified timestamp")
		return nil, nil
	}

	stats := make([]model.TickerStats, 0)
	for _, entry := range resp.ReadHistoricTickerStatsAtTimestamp {
		if len(entry.Stats) == 0 {
			log.Warn().Msgf("No TickerStats for timestamp: %d", entry.Timestamp)
			continue
		}
		for _, s := range entry.Stats {
			var liquidityEstimatePtr *string
			if s.LiquidityEstimate != "" {
				liquidityEstimatePtr = &s.LiquidityEstimate
			}
			stats = append(stats, model.TickerStats{
				Symbol:            s.Symbol,
				PriceChange:       s.PriceChange,
				PriceChangePct:    s.PriceChangePct,
				QuoteVolume:       s.QuoteVolume,
				Volume:            s.Volume,
				TradeCount:        s.TradeCount,
				HighPrice:         s.HighPrice,
				LowPrice:          s.LowPrice,
				LastPrice:         s.LastPrice,
				LiquidityEstimate: liquidityEstimatePtr,
			})
		}
	}
	return stats, nil
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
