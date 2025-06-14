package functions

import (
	"context"
	"strconv"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

// PairsOnTheMove is a function that takes in two slices of PriceData, currentPrices
// and previousPrices, and a float64 value marketMomentum, and returns a slice of
// Gainers structs and an error. The function calculates the percentage change
// between the current price and the previous price for each symbol in the currentPrices
// slice. If the percentage change is greater than or equal to marketMomentum, the
// function appends a Gainers struct to the PairsOnTheMoveList slice. The function
// uses the PercentageChange function to calculate the percentage change and
// strconv.ParseFloat function to parse the price data. The function prints out a
// message for each Pair that meets the marketMomentum condition.
func PairsOnTheMove(currentPrices, previousPrices []model.Pair, marketMomentum float64) (pairsOnTheMoveList []shared.Gainers, err error) {

	for i := 0; i < len(currentPrices); i++ {

		priceNow, _ := strconv.ParseFloat(currentPrices[i].Price, 64)
		priceNowSymbol := currentPrices[i].Symbol

		for i := range previousPrices {
			if previousPrices[i].Symbol == priceNowSymbol {
				previousPrice, _ := strconv.ParseFloat(previousPrices[i].Price, 64)
				change := shared.PercentageChange(previousPrice, priceNow)
				log.Debug().Float64("Price Now", priceNow).Float64("previous Price", previousPrice).Float64("change", change).Msg("Check")

				if change >= float64(marketMomentum) {
					pairsOnTheMoveList = append(pairsOnTheMoveList, shared.Gainers{Symbol: priceNowSymbol, IncrementPriceGain: change})
					log.Debug().Str("Current", currentPrices[0].Price).Str("Previous", previousPrices[0].Price).Str("Symbol:", priceNowSymbol).Float64("Increment Price Gain:", change).Msg("On the move")
				}
			}
		}
	}

	return pairsOnTheMoveList, nil
}

// FirstFilter is a function that takes in a path to historical price data,
// and a market momentum value. It then loads and compares the current and
// previous price data, and returns pairs of assets that have moved in the
// market by at least the specified momentum value. If an error occurs during
// any of the steps, an error is returned along with a nil slice of Gainers.
// Modified FirstFilter that uses in-memory enrichment instead of DB fetches
func FirstFilter(currentPrices []model.Pair, marketMomentum float64) ([]shared.Gainers, error) {
	var pairsOnTheMoveList []shared.Gainers

	for _, pair := range currentPrices {
		if pair.PercentageChange == nil {
			continue
		}
		change, err := strconv.ParseFloat(*pair.PercentageChange, 64)
		if err != nil {
			continue
		}
		if change >= marketMomentum {
			pairsOnTheMoveList = append(pairsOnTheMoveList, shared.Gainers{
				Symbol:             pair.Symbol,
				IncrementPriceGain: change,
			})
		}
	}
	return pairsOnTheMoveList, nil
}

// GetPriceData function retrieves price data from MongoDB using a GraphQL client
// and a timestamp. It returns a slice of PriceData.
func GetPriceData(ctx context.Context, client graphql.Client, datetime int, botName string) ([]model.Pair, error) {
	// Call the appropriate function to get historic prices at the specified timestamp

	log.Debug().Int("Time", datetime).Msg("getting data for")
	pricesList, err := graph.ReadHistoricPricesAtTimestamp(ctx, client, datetime)
	if err != nil {
		log.Error().Int("timestamp", datetime).Err(err).Msg("failed to load prices list from MongoDB")
		return nil, err
	}

	// Convert the GraphQL response to the desired output structure
	priceDataList, err := convertToPriceDataList(pricesList)
	if err != nil {
		log.Error().Int("timestamp", datetime).Err(err).Msg("failed to convert prices list")
		return nil, err
	}

	if len(priceDataList) == 0 {
		log.Warn().Int("Timestamp", datetime).Str("Bot", botName).Msg("No prices found for the specified timestamp")
	} else {
		log.Debug().Int("Timestamp", datetime).Interface("Prices", priceDataList).Msg("Prices fetched successfully")
	}

	return priceDataList, nil
}

func convertToPriceDataList(response *graph.ReadHistoricPricesAtTimestampResponse) ([]model.Pair, error) {
	var priceDataList []model.Pair

	for _, historicPrices := range response.GetReadHistoricPricesAtTimestamp() {
		for _, Pair := range historicPrices.GetPair() {
			priceData := model.Pair{
				Symbol: Pair.GetSymbol(),
				Price:  Pair.GetPrice(),
				// Add other fields as needed
			}
			priceDataList = append(priceDataList, priceData)
		}
	}

	return priceDataList, nil
}
