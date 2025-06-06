package functions

import (
	"context"
	"fmt"
	"sync"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	trade "cryptobotmanager.com/cbm-backend/microservices/binance/streamPrices"
	filter "cryptobotmanager.com/cbm-backend/microservices/filters/functions"
	reports "cryptobotmanager.com/cbm-backend/microservices/reports/functions"
	tradingBots "cryptobotmanager.com/cbm-backend/microservices/tradingBots/functions"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

func LetsTrade(ctx context.Context, client graphql.Client, market []model.Pair, datetime int) error {

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
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	log.Info().Int("Qty of pairs on the move", len(PairsOnTheMove)).Int("Datetime", datetime).Msg("FIRST FILTER (Is Active)- reduce all trading pairs to just those who have gained in the last 5 minutes")

	// Filter the current selection by weather is has shown enough liquidity to support a trade
	// After FirstFilter and before GetParameters
	log.Info().Msg("Filtering by liquidity...")
	LiquidPairsOnTheMove, err := filter.FilterByLiquidity(ctx, client, PairsOnTheMove, 500)
	if err != nil {
		log.Error().Err(err).Msg("Liquidity filter failed")
	}
	if len(LiquidPairsOnTheMove) == 0 {
		log.Warn().Msg("No liquid pairs found")
		return nil
	}
	log.Info().Int("Qty", len(LiquidPairsOnTheMove)).Msg("Pairs passing liquidity threshold")

	// Retrieve the details for all of the bots in the system
	log.Info().Msg("Loading strategy Details ...")
	strategyDetails, err := tradingBots.GetParameters(ctx, client)
	if err != nil {
		log.Error().Msgf("Failed to get strategy details!")
	}

	// Start a goroutine for each bot
	var wg sync.WaitGroup
	for _, details := range strategyDetails {
		wg.Add(1)
		go func(details model.StrategyInput) {
			defer wg.Done()

			//var chosenTicker string
			botName := details.BotInstanceName
			log.Info().Str("Name", botName).Int("Duration", details.TradeDuration).Int("IncrementsATR", details.IncrementsAtr).Int("ShortSMA", details.ShortSMADuration).Int("LongSMA", details.LongSMADuration).Msg("Strategy Details")

			fmt.Println("")
			log.Info().Msg("SECOND FILTER (SMA Gain)")
			coinsWithMomentum, err := filter.CompareSimpleMovingAverages(ctx, client, datetime, &LiquidPairsOnTheMove, details.ShortSMADuration, details.LongSMADuration, details.MovingAveMomentum, botName)
			if err != nil {
				log.Error().Msgf("coins With Momentum")
			}

			length := 0
			// Check if coinsWithMomentum is nil before using it
			if coinsWithMomentum != nil {
				length = len(*coinsWithMomentum)
				log.Info().Int("Qty", length).Msg("Coins with Market Momentum")

				if length == 0 {
					return
				}

				if length > 0 {
					fmt.Println("")
					log.Info().Msg("THIRD FILTER (Volatility)")
					log.Debug().Msg("Get Current Price And Calculate Average True Range")

					if len(*coinsWithMomentum) > 0 {
						chosenTicker := (*coinsWithMomentum)[0].Symbol

						log.Info().
							Str("Chosen Ticker (bypassed ATR)", chosenTicker).
							Msg("Paper Trading - Temporary Selection Without ATR")

						// Optionally call the MakeTrade or just log it for now
						trade.ListenAndPaperTrade(ctx, client, chosenTicker, details)
					}

					// for _, coin := range *coinsWithMomentum {
					// 	percentageGain, err := goBot.GetATR(tradeKlines, coin.Symbol, details.TradeDuration, details.IncrementsATR)
					// 	if err != nil {
					// 		log.Error().Msgf("Third filter percentage gain")
					// 	}

					// 	log.Debug().Str("Symbol", coin.Symbol).Float64("ATR Percentage Gain", percentageGain).Msg("Change")

					// 	for i := range *coinsWithMomentum {
					// 		if (*coinsWithMomentum)[i].Symbol == coin.Symbol {
					// 			(*coinsWithMomentum)[i].ATR = percentageGain
					// 			break
					// 		}
					// 	}
					// }

					// _, chosenTicker, _, err = goBot.FilterByAverageTrueRange(*coinsWithMomentum, details.MovingAveMomentum, cfg.WeightSMA, cfg.WeightATR, details.BotInstanceName)
					// if err != nil {
					// 	log.Error().Msgf("Filter By Average True Range!")
					// }

					// log.Info().Str("Chosen Ticker", chosenTicker).Msg("Paper Tading")
					// goBot.MakeTrade(client, chosenTicker, scenarioType, details)

				} else {
					log.Warn().Msg("coinsWithMomentum is nil")
					return
				}
			}
		}(details)
	}
	wg.Wait()
	return nil
}
