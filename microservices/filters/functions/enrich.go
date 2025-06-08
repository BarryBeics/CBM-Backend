package functions

import (
	"fmt"
	"strconv"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/rs/zerolog/log"
)

// enrichWithPercentageChange adds PercentageChange to each Pair by comparing with previousPrices
func EnrichWithPercentageChange(currentPrices, previousPrices []model.Pair) ([]model.Pair, error) {
	// Index previous prices for quick lookup
	prevMap := make(map[string]float64)
	for _, pair := range previousPrices {
		if price, err := strconv.ParseFloat(pair.Price, 64); err == nil {
			prevMap[pair.Symbol] = price
		}
	}

	// Enrich current prices with percentage change
	for i := range currentPrices {
		currentSymbol := currentPrices[i].Symbol
		currentPrice, err := strconv.ParseFloat(currentPrices[i].Price, 64)
		if err != nil {
			log.Warn().Str("symbol", currentSymbol).Str("price", currentPrices[i].Price).Msg("Invalid current price")
			continue
		}

		prevPrice, ok := prevMap[currentSymbol]
		if !ok {
			log.Warn().Str("symbol", currentSymbol).Msg("No matching previous price found")
			continue
		}

		change := shared.PercentageChange(prevPrice, currentPrice)
		percentStr := fmt.Sprintf("%.6f", change)
		currentPrices[i].PercentageChange = &percentStr

		log.Debug().
			Str("symbol", currentSymbol).
			Float64("prev", prevPrice).
			Float64("curr", currentPrice).
			Str("change", percentStr).
			Msg("Percentage change calculated")
	}

	return currentPrices, nil
}
