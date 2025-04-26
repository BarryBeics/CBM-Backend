package functions

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog/log"
)

// FetchPricesFromBinanceAPI fetches market prices from Binance API
// using API and Secret keys from environment variables.
// It returns a slice of PriceData structs and an error if any.
func FetchPricesFromBinanceAPI() (market []model.Pair, err error) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")
	client := binance.NewClient(apiKey, secretKey)
	if client == nil {
		return nil, fmt.Errorf("failed to create Binance client")
	}
	log.Info().Msg("Created conection to Binnance")

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

// SavePriceData writes market price data using GraphQL mutation
func SavePriceData(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {
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
	roundedEpochSeconds := shared.RoundTimeToFiveMinuteInterval(now) - 300

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
