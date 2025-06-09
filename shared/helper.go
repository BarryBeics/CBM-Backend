package shared

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

// RoundMinuteToFiveMinuteInterval function takes in a currentMinute integer as an
// argument, rounds it to the nearest five-minute interval and returns the rounded
// minute as an integer. It also logs debug messages using zerolog package.
// RoundTimeToFiveMinuteInterval takes in a time.Time value, rounds it to the nearest five-minute interval,
// and returns the rounded time.
func RoundTimeToFiveMinuteInterval(epochTime int64) int {
	// Convert epoch time to time.Time
	t := time.Unix(epochTime, 0)

	// Round to the nearest 5 minutes
	rounded := t.Round(5 * time.Minute)

	// Convert the rounded time back to epoch time
	int64Value := rounded.Unix()

	// Log rounded time
	//log.Trace().Int64("original_epoch_time", epochTime).Int64("rounded_epoch_time", int64Value).Msg("rounded time to five-minute interval")

	var roundedEpochTime int = int(int64Value)

	return roundedEpochTime
}

// CalculateMultiplier returns the multiplier corresponding to the given number.
// For example, when 'num' is 1, it returns 10; when 'num' is 2, it returns 100; and so on.
func CalculateMultiplier(num float64) float64 {
	if num == 1 {
		return 10
	} else if num == 2 {
		return 100
	} else if num == 3 {
		return 1000
	} else if num == 4 {
		return 10000
	} else {
		return 0
	}
}

// FindUniqueStrings returns the unique strings present in slice1 but not in slice2.
func FindUniqueStrings(slice1, slice2 []string) []string {
	unique := make(map[string]struct{})

	// Add items from slice1 to unique map
	for _, item := range slice1 {
		unique[item] = struct{}{}
	}

	// Remove items from slice2 that are already in the unique map
	for _, item := range slice2 {
		delete(unique, item)
	}

	// Convert the keys of the unique map to a slice
	uniqueItems := make([]string, 0, len(unique))
	for item := range unique {
		uniqueItems = append(uniqueItems, item)
	}

	return uniqueItems
}

// PercentageChange function takes in two float64 values as arguments, calculates
// the percentage change between the two values and returns it as a float64 value.
// It also logs debug messages using zerolog package.
func PercentageChange(inputOne, inputTwo float64) (result float64) {
	difference := inputTwo - inputOne
	result = difference / inputOne * 100

	// Log percentage change
	//log.Trace().Float64("input_one", inputOne).Float64("input_two", inputTwo).Float64("change", result).Msg("percentage change")

	return Round(result, 2)
}

func Round(x, decimal float64) float64 {
	unit := calculateValue(decimal)
	return math.Round(x*unit) / unit
}

func IncrementMean(m model.Mean, value float64) model.Mean {
	return model.Mean{
		Avg:   ((m.Avg * float64(m.Count)) + value) / float64(m.Count+1),
		Count: m.Count + 1,
	}
}

// RoundFloatToDecimal rounds the given floating-point number 'x'
// to the specified number of decimal places.
func RoundFloatToDecimal(x, decimal float64) float64 {
	if decimal < 1 || decimal > 4 {
		return x // Do not round if decimal is not 1, 2, or 3
	}
	unit := CalculateMultiplier(decimal)
	return math.Round(x*unit) / unit
}

// not sure what this is for - i will probably discover it when I rebuild the application
func calculateValue(num float64) float64 {
	if num == 1 {
		return 10
	} else if num == 2 {
		return 100
	} else if num == 3 {
		return 1000
	} else {
		return 0
	}
}

// // SavePriceData writes market price data using GraphQL mutation
// func SavePriceData(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {
// 	var pairsInput []graph.PairInput
// 	for _, price := range market {
// 		pairsInput = append(pairsInput, graph.PairInput{
// 			Symbol: price.Symbol,
// 			Price:  price.Price,
// 		})
// 	}

// 	chunks := chunkPairs(pairsInput, 250) // or 500 depending on how large each payload is
// 	for i, chunk := range chunks {
// 		log.Info().Int("chunk", i+1).Int("size", len(chunk)).Msg("Saving price chunk...")
// 		if err := saveChunkWithRetry(ctx, client, chunk, datetime, 3); err != nil {
// 			log.Error().Err(err).Int("chunk", i+1).Msg("Failed to save chunk after retries")
// 			return err // Or continue if partial failure is acceptable
// 		}
// 	}

// 	return nil
// }

// func SaveTradeStats(ctx context.Context, client graphql.Client, market []model.TickerStatsInput, datetime int) error {
// 	var chunks [][]model.TickerStatsInput
// 	chunkSize := 250 // or 500 depending on how large each payload is
// 	for chunkSize < len(market) {
// 		market, chunks = market[chunkSize:], append(chunks, market[0:chunkSize:chunkSize])
// 	}
// 	chunks = append(chunks, market)

// 	for i, chunk := range chunks {
// 		log.Info().Int("chunk", i+1).Int("size", len(chunk)).Msg("Saving trade stats chunk...")
// 		if err := saveTradeStatsChunkWithRetry(ctx, client, chunk, datetime, 3); err != nil {
// 			log.Error().Err(err).Int("chunk", i+1).Msg("Failed to save chunk after retries")
// 			return err // Or continue if partial failure is acceptable
// 		}
// 	}

// 	return nil
// }

// func chunkPairs(pairs []graph.PairInput, size int) [][]graph.PairInput {
// 	var chunks [][]graph.PairInput
// 	for size < len(pairs) {
// 		pairs, chunks = pairs[size:], append(chunks, pairs[0:size:size])
// 	}
// 	return append(chunks, pairs)
// }

// func saveChunkWithRetry(ctx context.Context, client graphql.Client, pairs []graph.PairInput, timestamp int, retries int) error {
// 	var err error
// 	for attempt := 1; attempt <= retries; attempt++ {
// 		input := graph.NewHistoricPriceInput{
// 			Pairs:     pairs,
// 			Timestamp: timestamp,
// 		}

// 		_, err = graph.CreateHistoricPrices(ctx, client, input)
// 		if err == nil {
// 			return nil
// 		}

// 		log.Warn().
// 			Err(err).
// 			Int("attempt", attempt).
// 			Int("chunk_size", len(pairs)).
// 			Msg("Failed to write price chunk, will retry...")

// 		time.Sleep(time.Duration(attempt*2) * time.Second) // backoff
// 	}
// 	return err
// }

// func saveTradeStatsChunkWithRetry(ctx context.Context, client graphql.Client, pairs []graph.TickerStatsInput, timestamp int, retries int) error {
// 	var err error
// 	for attempt := 1; attempt <= retries; attempt++ {
// 		input := graph.NewHistoricTickerStatsInput{
// 			Timestamp: timestamp,
// 			Stats:     pairs,
// 		}

// 		_, err = graph.CreateHistoricTickerStats(ctx, client, input)
// 		if err == nil {
// 			return nil
// 		}

// 		log.Warn().
// 			Err(err).
// 			Int("attempt", attempt).
// 			Int("chunk_size", len(pairs)).
// 			Msg("Failed to write trade stats chunk, will retry...")

// 		time.Sleep(time.Duration(attempt*2) * time.Second) // backoff
// 	}
// 	return err
// }

// SavePriceDataAsJSON saves Binance price data to a JSON file, appending data in 5-minute intervals
// func SavePriceDataAsJSON(market []model.Pair, timestamp int64) error {
// 	// Convert epoch timestamp to UTC date (YYYY-MM-DD)
// 	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

// 	// Define filename based on the date
// 	filename := fmt.Sprintf("/var/log/binance_prices_%s.json", date)

// 	// Read existing data (if file exists)
// 	var priceHistory []model.HistoricPrices
// 	file, err := os.Open(filename)
// 	if err == nil {
// 		defer file.Close()
// 		json.NewDecoder(file).Decode(&priceHistory)
// 	}

// 	// Convert []model.Pair to []*model.Pair
// 	var pairPtrs []*model.Pair
// 	for i := range market {
// 		pairPtrs = append(pairPtrs, &market[i])
// 	}

// 	// Append new data
// 	priceHistory = append(priceHistory, model.HistoricPrices{
// 		Pair:      pairPtrs,
// 		Timestamp: int(timestamp),
// 	})

// 	// Write updated JSON back to file
// 	fileData, err := json.MarshalIndent(priceHistory, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(filename, fileData, 0644)
// }

// // SavePriceDataAsJSON saves Binance price data to a JSON file, appending data in 5-minute intervals
// func SaveTradeStatsAsJSON(market []model.TickerStatsInput, timestamp int64) error {
// 	// Convert epoch timestamp to UTC date (YYYY-MM-DD)
// 	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

// 	// Define filename based on the date
// 	filename := fmt.Sprintf("/var/log/binance_prices_%s.json", date)

// 	// Read existing data (if file exists)
// 	var priceHistory []model.HistoricPrices
// 	file, err := os.Open(filename)
// 	if err == nil {
// 		defer file.Close()
// 		json.NewDecoder(file).Decode(&priceHistory)
// 	}

// 	// Convert []model.Pair to []*model.Pair
// 	var pairPtrs []*model.TickerStatsInput
// 	for i := range market {
// 		pairPtrs = append(pairPtrs, &market[i])
// 	}

// 	// Append new data
// 	priceHistory = append(priceHistory, model.HistoricPrices{
// 		Pair:      pairPtrs,
// 		Timestamp: int(timestamp),
// 	})

// 	// Write updated JSON back to file
// 	fileData, err := json.MarshalIndent(priceHistory, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(filename, fileData, 0644)
// }

// ParseDemoCSVData reads a CSV file and saves the data to the GraphQL database
func ParseDemoCSVData(ctx context.Context, client graphql.Client, csvFile string) ([]model.Pair, error) {
	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialize CSV reader
	reader := csv.NewReader(file)

	// Read the header row to get the symbol names
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Get Time
	now := time.Now().Unix()
	roundedEpochSeconds := RoundTimeToFiveMinuteInterval(now) - 300

	// Create the PriceData structure to hold the symbol-price pairs
	var currentPrices []model.Pair
	var rowCount int
	// Process each row and save data to the database
	for {
		rowCount++
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Create the PriceData structure to hold the symbol-price pairs
		var market []model.Pair

		// Process each column, starting from index 1
		for i := 0; i < len(headers); i++ {
			symbol := headers[i]
			price := record[i]

			// Create the PriceData struct for the symbol-price pair
			data := model.Pair{
				Symbol: symbol,
				Price:  price,
			}

			// Append the pair to the market slice
			market = append(market, data)

			// If i == 1, append the pair to the currentPrices slice
			if rowCount == 1 {
				currentPrices = append(currentPrices, data)
			}
		}

		// Save data to the database with the current roundedEpochSeconds
		err = SavePriceData(ctx, client, market, roundedEpochSeconds)
		if err != nil {
			log.Error().Err(err).Msgf("Save PriceData")
		}

		// Subtract 300 seconds (5 minutes) for the next row
		roundedEpochSeconds -= 300

		// Sleep for a short duration to avoid hitting rate limits (adjust as needed)
		time.Sleep(100 * time.Millisecond)
	}

	return currentPrices, nil
}

// NEW DRY CODE

func Chunk[T any](items []T, size int) [][]T {
	var chunks [][]T
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[:size:size])
	}
	return append(chunks, items)
}

func RetryWithBackoff[T any](ctx context.Context, data []T, retries int, fn func([]T) error) error {
	var err error
	for attempt := 1; attempt <= retries; attempt++ {
		if err = fn(data); err == nil {
			return nil
		}
		log.Warn().
			Err(err).
			Int("attempt", attempt).
			Int("chunk_size", len(data)).
			Msg("Failed to write chunk, will retry...")
		time.Sleep(time.Duration(attempt*2) * time.Second)
	}
	return err
}

func SaveWithChunks[T any](ctx context.Context, client graphql.Client, items []T, chunkSize int, retries int, mutationFn func(context.Context, graphql.Client, []T, int) error, timestamp int) error {
	chunks := Chunk(items, chunkSize)
	for i, chunk := range chunks {
		log.Debug().Int("chunk", i+1).Int("size", len(chunk)).Msg("Saving chunk...")
		err := RetryWithBackoff(ctx, chunk, retries, func(data []T) error {
			return mutationFn(ctx, client, data, timestamp)
		})
		if err != nil {
			log.Error().Err(err).Int("chunk", i+1).Msg("Failed to save chunk after retries")
			return err
		}
	}
	return nil
}

func safeDereference(s *string, fallback string) string {
	if s != nil {
		return *s
	}
	return fallback
}

func SavePriceData(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {
	pairsInput := make([]graph.PairInput, 0, len(market))
	for _, p := range market {
		pairsInput = append(pairsInput, graph.PairInput{
			Symbol:           p.Symbol,
			Price:            p.Price,
			PercentageChange: safeDereference(p.PercentageChange, "0"),
		})
		// if p.PercentageChange == nil {
		// 	log.Warn().Str("symbol", p.Symbol).Msg("Missing PercentageChange")
		// }

	}

	return SaveWithChunks(ctx, client, pairsInput, 250, 3, func(ctx context.Context, client graphql.Client, chunk []graph.PairInput, timestamp int) error {
		input := graph.NewHistoricPriceInput{
			Pairs:     chunk,
			Timestamp: timestamp,
		}
		_, err := graph.CreateHistoricPrices(ctx, client, input)
		return err
	}, datetime)
}

func SaveTradeStats(ctx context.Context, client graphql.Client, market []model.TickerStatsInput, datetime int) error {
	statsInput := make([]graph.TickerStatsInput, 0, len(market))
	for _, s := range market {
		var liquidityEstimate string
		if s.LiquidityEstimate != nil {
			liquidityEstimate = *s.LiquidityEstimate
		}
		statsInput = append(statsInput, graph.TickerStatsInput{
			Symbol:            s.Symbol,
			PriceChange:       s.PriceChange,
			PriceChangePct:    s.PriceChangePct,
			QuoteVolume:       s.QuoteVolume,
			Volume:            s.Volume,
			TradeCount:        s.TradeCount,
			HighPrice:         s.HighPrice,
			LowPrice:          s.LowPrice,
			LastPrice:         s.LastPrice,
			LiquidityEstimate: liquidityEstimate,
		})
	}

	return SaveWithChunks(ctx, client, statsInput, 250, 3, func(ctx context.Context, client graphql.Client, stats []graph.TickerStatsInput, ts int) error {
		_, err := graph.CreateHistoricTickerStats(ctx, client, graph.NewHistoricTickerStatsInput{
			Timestamp: ts,
			Stats:     stats,
		})
		return err
	}, datetime)
}

func saveJSON[TInput any, TBatch any](
	filenameFormat string,
	market []TInput,
	timestamp int64,
	toBatch func([]TInput, int) TBatch,
) error {
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	filename := fmt.Sprintf(filenameFormat, date)

	// Read existing data
	var history []TBatch
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		json.NewDecoder(file).Decode(&history)
	}

	// Append new data
	history = append(history, toBatch(market, int(timestamp)))

	// Write back
	fileData, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, fileData, 0644)
}

func SavePriceDataAsJSON(market []model.Pair, timestamp int64) error {
	return saveJSON(
		"/var/log/binance_prices_%s.json",
		market,
		timestamp,
		func(pairs []model.Pair, ts int) model.HistoricPrices {
			ptrs := make([]*model.Pair, len(pairs))
			for i := range pairs {
				ptrs[i] = &pairs[i]
			}
			return model.HistoricPrices{
				Pair:      ptrs,
				Timestamp: ts,
			}
		},
	)
}

func SaveTradeStatsAsJSON(market []model.TickerStatsInput, timestamp int64) error {
	return saveJSON(
		"/var/log/binance_stats_%s.json",
		market,
		timestamp,
		func(inputs []model.TickerStatsInput, ts int) model.HistoricTickerStats {
			stats := make([]*model.TickerStats, len(inputs))
			for i, in := range inputs {
				stats[i] = &model.TickerStats{
					Symbol:            in.Symbol,
					PriceChange:       in.PriceChange,
					PriceChangePct:    in.PriceChangePct,
					QuoteVolume:       in.QuoteVolume,
					Volume:            in.Volume,
					TradeCount:        in.TradeCount,
					HighPrice:         in.HighPrice,
					LowPrice:          in.LowPrice,
					LastPrice:         in.LastPrice,
					LiquidityEstimate: in.LiquidityEstimate,
				}
			}
			return model.HistoricTickerStats{
				Stats:     stats,
				Timestamp: ts,
				CreatedAt: time.Now().UTC(),
			}
		},
	)
}

func GetPreviousTime(currentTime int, minutesToSubtract int) (int64, error) {
	// Convert the integer currentTime to time.Time
	currentTimeAsTime := time.Unix(int64(currentTime), 0)

	// Subtract minutes from the current time
	subtractedTime := currentTimeAsTime.Add(time.Duration(-minutesToSubtract) * time.Minute)

	// Convert the result to epoch time (seconds since 1970)
	epochTime := subtractedTime.Unix()

	return epochTime, nil
}
