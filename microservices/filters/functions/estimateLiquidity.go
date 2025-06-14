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

	log.Info().
		Int("pair_count", len(pairs)).
		Float64("threshold", threshold).
		Msg("Starting liquidity filter")

	for _, pair := range pairs {
		log.Debug().
			Str("symbol", pair.Symbol).
			Msg("Fetching ticker stats")

		stats, err := graph.ReadTickerStatsBySymbol(ctx, client, pair.Symbol, 1)
		if err != nil {
			log.Error().
				Str("symbol", pair.Symbol).
				Err(err).
				Msg("Failed to load liquidity stats")
			continue
		}

		thing := stats.ReadTickerStatsBySymbol
		if len(thing) == 0 {
			log.Warn().
				Str("symbol", pair.Symbol).
				Msg("No ticker stats returned")
			continue
		}

		rawLiquidity := thing[0].LiquidityEstimate
		if rawLiquidity == "" {
			log.Warn().
				Str("symbol", pair.Symbol).
				Msg("LiquidityEstimate is empty")
			continue
		}

		liquidity, err := strconv.ParseFloat(rawLiquidity, 64)
		if err != nil {
			log.Error().
				Str("symbol", pair.Symbol).
				Str("liquidity", rawLiquidity).
				Err(err).
				Msg("Failed to parse LiquidityEstimate")
			continue
		}

		log.Debug().
			Str("symbol", pair.Symbol).
			Str("raw", rawLiquidity).
			Float64("parsed", liquidity).
			Msg("Parsed liquidity successfully")

		if liquidity >= threshold {
			log.Info().
				Str("symbol", pair.Symbol).
				Float64("liquidity", liquidity).
				Msg("Pair meets liquidity threshold")
			liquidPairs = append(liquidPairs, pair)
		} else {
			log.Debug().
				Str("symbol", pair.Symbol).
				Float64("liquidity", liquidity).
				Msg("Pair below liquidity threshold")
		}
	}

	log.Info().
		Int("qualified_pairs", len(liquidPairs)).
		Msg("Completed liquidity filtering")

	return liquidPairs, nil
}
