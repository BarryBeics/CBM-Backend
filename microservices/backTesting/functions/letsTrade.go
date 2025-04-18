package functions

import (
	"context"
	"fmt"

	filter "cryptobotmanager.com/cbm-backend/microservices/filters/functions"
	reports "cryptobotmanager.com/cbm-backend/microservices/reports/functions"
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
	PairsOnTheMove, err := filter.FirstFilter(ctx, client, datetime, cfg.ActiveMarketThreshold)
	if err != nil {
		log.Error().Msgf("Pairs on the move!")
	}

	if len(PairsOnTheMove) == 0 {
		return nil
	}
	reports.MarketActivityReport(client, cfg.TopAverages, PairsOnTheMove, datetime)
	log.Info().Int("Qty of pairs on the move", len(PairsOnTheMove)).Msg("FIRST FILTER (Is Active)- reduce all trading pairs to just those who have gained in the last 5 minutes")
	fmt.Println("")

	return nil
}
