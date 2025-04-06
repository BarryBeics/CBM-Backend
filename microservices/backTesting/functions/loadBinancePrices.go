package functions

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type PriceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// SavePriceData writes market price data using GraphQL mutation
func SavePriceData(ctx context.Context, client graphql.Client, market []PriceData, datetime int) error {
	// Create an array of PairInput from PriceData
	var pairsInput []graph.PairInput
	for _, price := range market {
		pairInput := graph.PairInput{
			Symbol: price.Symbol,
			Price:  price.Price,
		}

		pairsInput = append(pairsInput, pairInput)

	}

	// Create NewHistoricPriceInput
	input := graph.NewHistoricPriceInput{
		Pairs:     pairsInput,
		Timestamp: datetime,
	}

	// Call the GraphQL mutation
	_, err := graph.CreateHistoricPrices(ctx, client, input)
	if err != nil {
		log.Printf("Error creating historic prices: %v", err)
		return err
	}

	return nil
}

type ExportPriceData struct {
	Pairs     []PriceData `json:"pairs"`
	Timestamp int64       `json:"timestamp"`
}

// ParseDemoCSVData reads a CSV file and saves the data to the GraphQL database
func ParseDemoCSVData(ctx context.Context, client graphql.Client, csvFile string) ([]PriceData, error) {
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
	roundedEpochSeconds := shared.RoundTimeToFiveMinuteInterval(now) - 300

	// Create the PriceData structure to hold the symbol-price pairs
	var currentPrices []PriceData
	var rowCount int
	// Process each row and save data to the database
	for {
		rowCount++
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Create the PriceData structure to hold the symbol-price pairs
		var market []PriceData

		// Process each column, starting from index 1
		for i := 0; i < len(headers); i++ {
			symbol := headers[i]
			price := record[i]

			// Create the PriceData struct for the symbol-price pair
			data := PriceData{
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

func ExtractTimestampFromFilename(filename string) int {
	base := filepath.Base(filename)
	parts := strings.TrimSuffix(strings.TrimPrefix(base, "binance_prices_"), ".json")

	// Parse date like 2025-03-30
	t, err := time.Parse("2006-01-02", parts)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to parse date from filename")
		return int(time.Now().Unix())
	}
	// Add a fake time (e.g. noon)
	t = t.Add(12 * time.Hour)
	return int(t.Unix())
}

func LoadPriceSnapshotsFromFile(path string) ([]ExportPriceData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var snapshots []ExportPriceData
	err = json.Unmarshal(bytes, &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
