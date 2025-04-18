package functions

import (
	"context"

	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type MutationTradeOutcomeReport struct {
	CreateTradeOutcomeReport struct {
		ID               string  `json:"_id"`
		Timestamp        int     `json:"Timestamp"`
		BotName          string  `json:"BotName"`
		PercentageChange float64 `json:"PercentageChange"`
		Balance          float64 `json:"Balance"`
		Symbol           string  `json:"Symbol"`
		Outcome          string  `json:"Outcome"`
		ElapsedTime      int     `json:"ElapsedTime"`
		Volume           float64 `json:"Volume"`
		FearGreedIndex   int     `json:"FearGreedIndex"`
		MarketStatus     string  `json:"MarketStatus"`
	} `json:"createTradeOutcomeReport"`
}

type NewTradeOutcomeReport struct {
	Timestamp        int     `json:"Timestamp"`
	BotName          string  `json:"BotName"`
	PercentageChange float64 `json:"PercentageChange"`
	Balance          float64 `json:"Balance"`
	Symbol           string  `json:"Symbol"`
	Outcome          string  `json:"Outcome"`
	ElapsedTime      int     `json:"ElapsedTime"`
	Volume           float64 `json:"Volume"`
	FearGreedIndex   int     `json:"FearGreedIndex"`
	MarketStatus     string  `json:"MarketStatus"`
}

func TradeOutcomeReport(client graphql.Client, timeStamp, elapsedTime int, botName string, PercentageChange, updatedBalance, volume, Fee float64, symbol, outcome string) {

	ctx := context.Background()

	// Make the GraphQL mutation request
	_, err := graph.CreateTradeOutcomeReport(ctx, client,
		timeStamp,
		botName,
		PercentageChange,
		updatedBalance,
		symbol,
		outcome,
		Fee,
		int(elapsedTime),
		volume,
		70,
		"TBC",
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to add trade outcome report")
	}

}
