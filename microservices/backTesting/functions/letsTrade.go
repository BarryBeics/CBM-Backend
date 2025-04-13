package functions

import (
	"context"
	"fmt"

	filterFunctions "cryptobotmanager.com/cbm-backend/microservices/filters/functions"
	"cryptobotmanager.com/cbm-backend/resolvers/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func letTrade(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {

	err := SavePriceData(ctx, client, market, datetime)
	if err != nil {
		log.Error().Err(err).Int("timestamp", datetime).Msg("Failed to save snapshot")
	} else {
		log.Info().Int("timestamp", datetime).Msg("Replayed snapshot")
	}

	cfg := shared.GetDefaultCfg()

	// Report on the market activity
	PairsOnTheMove, err := filterFunctions.FirstFilter(ctx, client, datetime, cfg.ActiveMarketThreshold)
	if err != nil {
		log.Error().Msgf("Pairs on the move!")
	}

	fmt.Print(PairsOnTheMove)

	return nil
}
