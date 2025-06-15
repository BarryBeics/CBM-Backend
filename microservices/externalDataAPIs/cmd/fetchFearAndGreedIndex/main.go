package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type FNGResponse struct {
	Data []struct {
		Value               string `json:"value"`
		ValueClassification string `json:"value_classification"`
		Timestamp           string `json:"timestamp"`
	} `json:"data"`
}

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

	// Step 1: Fetch from FNG API
	resp, err := http.Get("https://api.alternative.me/fng/?limit=1")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch FNG index")
	}
	defer resp.Body.Close()

	var raw FNGResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode response")
	}

	if len(raw.Data) == 0 {
		log.Error().Msg("No data returned from FNG API")
		return
	}

	latest := raw.Data[0]

	tsInt, err := strconv.Atoi(latest.Timestamp)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid timestamp")
	}

	fmt.Printf("Fear & Greed Index: %s (%s) at %d\n", latest.Value, latest.ValueClassification, tsInt)

	// Step 3: Prepare GraphQL client and input
	client := graphql.NewClient(backend, &http.Client{})
	ctx := context.Background()

	req, err := graph.UpsertFearAndGreedIndex(ctx, client, tsInt, latest.Value, latest.ValueClassification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GraphQL request")
	}

	fmt.Printf("Saved to DB: %+v\n", req.UpsertFearAndGreedIndex)
}
