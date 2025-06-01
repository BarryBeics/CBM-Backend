package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"cryptobotmanager.com/cbm-backend/microservices/backTesting/functions"
	"cryptobotmanager.com/cbm-backend/shared"
	sharedlog "github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: No .env file found or failed to load")
	}

	// Initialize logger
	shared.SetupLogger()

	if len(flag.Args()) == 1 {
		switch flag.Arg(0) {
		case "Report":
			functions.PrintFunctionsWithoutTestCoverage()
			return
		default:
			fmt.Println("Invalid argument. Usage: go run main.go [Report]")
			os.Exit(1)
		}
	}

	// Setup GraphQL backend client
	backend := os.Getenv("TRADING_BOT_URL")
	if backend == "" {
		backend = "http://cbm-api:8080/query"
	}

	fmt.Println("SYSTEM_MODE is:", os.Getenv("SYSTEM_MODE"))

	// Setup GraphQL backend client
	// if os.Getenv("SYSTEM_MODE") == "local" {
	err = functions.CSVPrices(backend)
	if err != nil {
		sharedlog.Error().Err(err).Msg("Failed to load JSON data")
	}

	// } else {
	// 	err := functions.BinancePrices(backend)
	// 	if err != nil {
	// 		sharedlog.Error().Err(err).Msgf("Failed to get price data from Binance!")
	// 	}
	// }

}

// func BinancePrices(backend string) error {
// 	// make live API calls to Binance
// 	// Get the nearest whole 5 minutes & print the current time
// 	now := time.Now().Unix()
// 	roundedEpoch := shared.RoundTimeToFiveMinuteInterval(now)
// 	log.Info().Int64("Executing task at:", now).Int("Rounded time", roundedEpoch).Msg("Time")

// 	// Create Client & Context
// 	client := graphql.NewClient(backend, &http.Client{})
// 	ctx := context.Background()
// 	var market []model.Pair
// 	var err error

// 	market, err = binance.FetchPricesFromBinanceAPI()
// 	if err != nil {
// 		log.Error().Err(err).Msgf("Failed to get price data from Binance!")
// 	}

// 	err = shared.SavePriceDataAsJSON(market, int64(roundedEpoch))
// 	if err != nil {
// 		log.Error().Err(err).Msgf("Failed to save price data to JSON!")
// 	}

// 	if err := shared.SavePriceData(ctx, client, market, roundedEpoch); err != nil {
// 		log.Error().Err(err).Int("timestamp", roundedEpoch).Msg("Save PriceData")
// 	}
// 	return nil
// }
