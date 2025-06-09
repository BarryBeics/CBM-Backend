package functions

import (
	"context"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func MarketActivityReport(client graphql.Client, TopAverages []int, pairsOnTheMove []shared.Gainers, now int) {
	if len(pairsOnTheMove) == 0 {
		log.Warn().Msg("No pairs on the move, skipping market activity report")
		return
	}

	// Log the number of pairs on the move
	log.Debug().Int("pairs_count", len(pairsOnTheMove)).Msg("the market activity report")

	for _, pair := range pairsOnTheMove {
		log.Info().Str("Symbol", pair.Symbol).
			Float64("PriceGain", pair.SMAPriceGain).
			Msg("Pair on the move")
	}

	ManageSymbolStats(client, pairsOnTheMove[:min(10, len(pairsOnTheMove))])

	// Call the ActivityReport function to handle the report creation
	ActivityReport(client, TopAverages, pairsOnTheMove, now)

}

func ActivityReport(client graphql.Client, TopAverages []int, pairsOnTheMove []shared.Gainers, now int) {
	allPairs := len(pairsOnTheMove)
	allMovers := AverageGain(&pairsOnTheMove, allPairs)

	// Log the number of pairs on the move
	log.Debug().Int("pairs_count", allPairs).Msg("the market activity report")
	ctx := context.Background()

	const fearGreedIndex = 80 // Set a constant value for fearGreedIndex

	if allPairs != 0 {
		// Calculate topN averages
		var topA, topB, topC float64
		if len(pairsOnTheMove) >= TopAverages[0] {
			topA = AverageGain(&pairsOnTheMove, TopAverages[0])
		}
		if len(pairsOnTheMove) >= TopAverages[1] {
			topB = AverageGain(&pairsOnTheMove, TopAverages[1])
		}
		if len(pairsOnTheMove) >= TopAverages[2] {
			topC = AverageGain(&pairsOnTheMove, TopAverages[2])
		}

		// Make the GraphQL mutation request for each numberOfCoins
		for _, numberOfCoins := range TopAverages {
			if numberOfCoins+1 > allPairs {
				break
			}
			avgGain := AverageGain(&pairsOnTheMove, numberOfCoins)

			log.Debug().Int("time", now).Int("Qty", numberOfCoins).Float64("avgGain", avgGain).Msg("Check")
		}

		// Move the GraphQL mutation request here with the correct values for topA, topB, topC
		_, err := graph.CreateActivityReport(ctx, client,
			now,
			allPairs,
			allMovers,
			topA,
			topB,
			topC,
			fearGreedIndex)

		if err != nil {
			log.Error().Err(err).Msg("failed to add activity report")
		}
	}
}

func ManageSymbolStats(client graphql.Client, top10 []shared.Gainers) {
	ctx := context.Background()

	for i, pair := range top10 {
		// Fetch existing stats (if any)
		res, err := graph.SymbolStatsBySymbol(ctx, client, pair.Symbol)
		if err != nil {
			log.Warn().Str("symbol", pair.Symbol).Msg("No existing stats found, creating new entry")

			// Create default PositionCounts with zeros
			positionCounts := make([]int, 10)
			if i < 10 {
				positionCounts[i] = 1
			}

			// Create new SymbolStats entry
			_, err := graph.CreateSymbolStats(ctx, client, graph.CreateSymbolStatsInput{
				Symbol:         pair.Symbol,
				PositionCounts: positionCounts,
			})
			if err != nil {
				log.Error().Err(err).Str("symbol", pair.Symbol).Msg("Failed to create SymbolStats")
			}
			continue
		}

		// Update PositionCounts
		newCounts := append([]int(nil), res.PositionCounts...) // copy to avoid mutation
		if i < len(newCounts) {
			newCounts[i]++
		} else {
			// Extend the slice if needed (shouldn't happen if fixed at 10)
			for len(newCounts) <= i {
				newCounts = append(newCounts, 0)
			}
			newCounts[i] = 1
		}

		// Upsert
		_, err = graph.UpsertSymbolStats(ctx, client, model.UpsertSymbolStatsInput{
			Symbol:               pair.Symbol,
			PositionCounts:       newCounts,
			AvgLiquidityEstimate: newAvg,
			MaxLiquidityEstimate: newMax,
			MinLiquidityEstimate: newMin,
		})
		if err != nil {
			log.Error().Err(err).Str("symbol", pair.Symbol).Msg("Failed to update SymbolStats")
		}
	}
}
