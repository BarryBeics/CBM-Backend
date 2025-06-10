package functions

import (
	"context"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	sharedmodel "cryptobotmanager.com/cbm-backend/shared/model"
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
		res, err := graph.GetSymbolStatsBySymbol(ctx, client, pair.Symbol)
		// Always create a fresh slice for positionCounts
		positionCounts := make([]*model.MeanInput, 10)
		for j := 0; j < 10; j++ {
			positionCounts[j] = &model.MeanInput{Avg: 0, Count: 0}
		}

		if err != nil || res == nil || res.SymbolStatsBySymbol.Symbol == "" {
			log.Warn().Str("symbol", pair.Symbol).Msg("No existing stats found, creating new entry")
			// New entry: set this position
			if i < 10 {
				positionCounts[i] = &model.MeanInput{Avg: pair.IncrementPriceGain, Count: 1}
			}
		} else {
			// Existing stats: copy over old values and update this position
			for j := 0; j < len(res.SymbolStatsBySymbol.PositionCounts) && j < 10; j++ {
				old := res.SymbolStatsBySymbol.PositionCounts[j]
				positionCounts[j] = &model.MeanInput{Avg: old.Avg, Count: old.Count}
			}
			if i < 10 {
				oldStat := res.SymbolStatsBySymbol.PositionCounts[i]
				newCount := oldStat.Count + 1
				newAvg := ((oldStat.Avg * float64(oldStat.Count)) + pair.IncrementPriceGain) / float64(newCount)
				positionCounts[i] = &model.MeanInput{
					Avg:   newAvg,
					Count: newCount,
				}
			}
		}

		// Convert to []model.MeanInput just before passing to UpsertPositionCounts
		sharedModelPositionCounts := make([]sharedmodel.MeanInput, 0, len(positionCounts))
		for _, v := range positionCounts {
			sharedModelPositionCounts = append(sharedModelPositionCounts, sharedmodel.MeanInput{
				Avg:   v.Avg,
				Count: v.Count,
			})
		}

		_, err = graph.UpsertPositionCounts(ctx, client, pair.Symbol, sharedModelPositionCounts)
		if err != nil {
			log.Error().Err(err).Str("symbol", pair.Symbol).Msg("Failed to create or update SymbolStats")
		}
	}
}
