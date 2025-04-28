package main

// Add any other necessary imports

func main() {

	// backend := os.Getenv("TRADING_BOT_URL")
	// if backend == "" {
	// 	backend = "http://cbm-api:8080/query"
	// }

	// // Get the nearest whole 5 minutes & print the current time
	// now := time.Now().Unix()
	// roundedEpochSeconds := shared.RoundTimeToFiveMinuteInterval(now)
	// log.Info().Int64("Executing task at:", now).Int("Rounded time", roundedEpochSeconds).Msg("Time")

	// // Create Client & Context
	// client := graphql.NewClient(backend, &http.Client{})
	// ctx := context.Background()
	// var market []model.Pair
	// var err error

	// market, err = functions.FetchPricesFromBinanceAPI()
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Failed to get price data from Binance!")
	// }

	// err = functions.SavePriceDataAsJSON(market, int64(roundedEpochSeconds))
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Failed to save price data to JSON!")
	// }

	// err = functions.SavePriceData(ctx, client, market, roundedEpochSeconds)
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Save PriceData")
	// }

}
