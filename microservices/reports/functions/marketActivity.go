package functions

import (
	"context"

	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type MutationActivityReport struct {
	CreateActivityReport struct {
		ID             string  `json:"_id"`
		Timestamp      int     `json:"Timestamp"`
		Qty            int     `json:"Qty"`
		AvgGain        float64 `json:"AvgGain"`
		TopAGain       float64 `json:"TopAGain"`
		TopBGain       float64 `json:"TopBGain"`
		TopCGain       float64 `json:"TopCGain"`
		FearGreedIndex int     `json:"FearGreedIndex"`
	} `json:"createActivityReport"`
}

type NewActivityReport struct {
	Timestamp      int     `json:"Timestamp"`
	Qty            int     `json:"Qty"`
	AvgGain        float64 `json:"AvgGain"`
	TopAGain       float64 `json:"TopAGain"`
	TopBGain       float64 `json:"TopBGain"`
	TopCGain       float64 `json:"TopCGain"`
	FearGreedIndex int     `json:"FearGreedIndex"`
}

func MarketActivityReport(client graphql.Client, TopAverages []int, pairsOnTheMove []shared.Gainers, now int) {
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
