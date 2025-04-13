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
		backend = "http://resolvers:8080/query"
	}

	fmt.Println("SYSTEM_MODE is:", os.Getenv("SYSTEM_MODE"))

	// Setup GraphQL backend client
	if os.Getenv("SYSTEM_MODE") == "local" {
		err := functions.CSVPrices(backend)
		if err != nil {
			sharedlog.Error().Err(err).Msg("Failed to load JSON data")
		}

	} else {
		err := functions.BinancePrices(backend)
		if err != nil {
			sharedlog.Error().Err(err).Msgf("Failed to get price data from Binance!")
		}
	}

}
