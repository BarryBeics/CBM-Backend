package functions

import (
	"context"

	"strconv"

	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func FilterByLiquidity(ctx context.Context, client graphql.Client, pairs []shared.Gainers, threshold float64) ([]shared.Gainers, error) {
	var liquidPairs []shared.Gainers

	log.Info().Int("qty pairs", len(pairs)).Msg("filtering pairs by liquidity")

	for _, pair := range pairs {
		stats, err := graph.GetTickerLiquidityEstimate(ctx, client, pair.Symbol, 1)
		if err != nil {
			log.Error().Str("symbol", pair.Symbol).Err(err).Msg("failed to load liquidity stats")
			continue // Don't return early, just skip this one
		}

		thing := stats.GetTickerStatsBySymbol

		if len(thing) == 0 || thing[0].LiquidityEstimate == "" {
			continue
		}

		liquidityStr := thing[0].LiquidityEstimate
		liquidity, err := strconv.ParseFloat(liquidityStr, 64)
		if err != nil {
			log.Error().Str("symbol", pair.Symbol).Str("liquidity", liquidityStr).Err(err).Msg("failed to parse liquidity estimate")
			continue
		}
		if liquidity >= threshold {
			liquidPairs = append(liquidPairs, pair)
			log.Info().Str("symbol", pair.Symbol).
				Float64("liquidity", liquidity).
				Msg("pair meets liquidity threshold")
		}

	}

	return liquidPairs, nil
}
