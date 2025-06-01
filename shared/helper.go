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

// SavePriceData writes market price data using GraphQL mutation
func SavePriceData(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {
	var pairsInput []graph.PairInput
	for _, price := range market {
		pairsInput = append(pairsInput, graph.PairInput{
			Symbol: price.Symbol,
			Price:  price.Price,
		})
	}

	chunks := chunkPairs(pairsInput, 250) // or 500 depending on how large each payload is
	for i, chunk := range chunks {
		log.Info().Int("chunk", i+1).Int("size", len(chunk)).Msg("Saving price chunk...")
		if err := saveChunkWithRetry(ctx, client, chunk, datetime, 3); err != nil {
			log.Error().Err(err).Int("chunk", i+1).Msg("Failed to save chunk after retries")
			return err // Or continue if partial failure is acceptable
		}
	}

	return nil
}

func chunkPairs(pairs []graph.PairInput, size int) [][]graph.PairInput {
	var chunks [][]graph.PairInput
	for size < len(pairs) {
		pairs, chunks = pairs[size:], append(chunks, pairs[0:size:size])
	}
	return append(chunks, pairs)
}

func saveChunkWithRetry(ctx context.Context, client graphql.Client, pairs []graph.PairInput, timestamp int, retries int) error {
	var err error
	for attempt := 1; attempt <= retries; attempt++ {
		input := graph.NewHistoricPriceInput{
			Pairs:     pairs,
			Timestamp: timestamp,
		}

		_, err = graph.CreateHistoricPrices(ctx, client, input)
		if err == nil {
			return nil
		}

		log.Warn().
			Err(err).
			Int("attempt", attempt).
			Int("chunk_size", len(pairs)).
			Msg("Failed to write price chunk, will retry...")

		time.Sleep(time.Duration(attempt*2) * time.Second) // backoff
	}
	return err
}

// SavePriceDataAsJSON saves Binance price data to a JSON file, appending data in 5-minute intervals
func SavePriceDataAsJSON(market []model.Pair, timestamp int64) error {
	// Convert epoch timestamp to UTC date (YYYY-MM-DD)
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	// Define filename based on the date
	filename := fmt.Sprintf("/var/log/binance_prices_%s.json", date)

	// Read existing data (if file exists)
	var priceHistory []model.HistoricPrices
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		json.NewDecoder(file).Decode(&priceHistory)
	}

	// Convert []model.Pair to []*model.Pair
	var pairPtrs []*model.Pair
	for i := range market {
		pairPtrs = append(pairPtrs, &market[i])
	}

	// Append new data
	priceHistory = append(priceHistory, model.HistoricPrices{
		Pair:      pairPtrs,
		Timestamp: int(timestamp),
	})

	// Write updated JSON back to file
	fileData, err := json.MarshalIndent(priceHistory, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, fileData, 0644)
}

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
