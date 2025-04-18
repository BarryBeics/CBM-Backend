package functions

import (
	"context"
	"fmt"
	"strconv"

	"cryptobotmanager.com/cbm-backend/resolvers/graph/model"
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type SMA struct {
	TimeFrame      int
	PriceDataArray []model.Pair
}

type RollingTotal struct {
	Symbol       string
	RollingTotal float64
}

// CompareSimpleMovingAverages iterates through the list of crypto pairs to determine if their simple moving averages (SMAs)
// for a given short and long period indicate momentum in the market. It calculates the SMAs for each time frame of historic
// data loaded from JSON files, and compares the short and long SMAs. If the short SMA is greater than the long SMA, it
// calculates the percentage gain and adds the coin symbol and gain to the coinsWithMomentum slice if it meets or exceeds
// the movingAveMomentum threshold. If no coins meet the threshold, a warning message is logged.
func CompareSimpleMovingAverages(ctx context.Context, client graphql.Client, datetime int, trackGainers *[]shared.Gainers, short, long int, movingAveMomentum float64, botName string) (*[]shared.Gainers, error) {

	log.Debug().Msg("Extracting price data from database ...")
	allPriceData, _ := extractPriceData(ctx, client, trackGainers, datetime, long, botName)

	fmt.Print(allPriceData)
	fmt.Print("INSIDE COmpare SMA")

	currentPrices, shortAverages, longAverages, err := ProcessAllPriceData(allPriceData, short, long, botName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	log.Debug().Int("current Prices", len(currentPrices)).Int("Short Averages", len(shortAverages)).Int("Long Averages", len(longAverages)).Msg("Counts")

	log.Debug().Msg("Comparing short period SMAverages against long period SMAverages")
	// Compare shortAverages against longAverages
	for symbol, shortAvg := range shortAverages {
		if longAvg, found := longAverages[symbol]; found {
			log.Debug().Str("Symbol", symbol).Float64("ShortAvg", shortAvg).Float64("LongAvg", longAvg).Msg("Comparing shortAvg and longAvg")
			if shortAvg > longAvg {
				log.Debug().Str("Symbol", symbol).Msg("ShortAvg is greater than LongAvg")
				if _, exists := currentPrices[symbol]; exists {
					percentageGain := shared.PercentageChange(longAvg, shortAvg)
					log.Debug().Float64("longAve", longAvg).Float64("short ave", shortAvg).Float64("Percantage Gain", percentageGain).Str("Symbol", symbol).Msg("Moving into short term market momentum")

					// Find the Gainers struct in trackGainers slice and update its fields if the percentage gain is greater than the target moving ave momentum
					if percentageGain > movingAveMomentum {
						for i := range *trackGainers {
							if (*trackGainers)[i].Symbol == symbol {
								(*trackGainers)[i].SMAPriceGain = percentageGain // Update the SMAPriceGain field
								break
							}
						}
					}
				} else {
					fmt.Printf("Current price not found for symbol %s\n", symbol)
				}
			}
		} else {
			log.Debug().Str("Symbol", symbol).Msg("LongAvg not found for symbol")
		}
	}

	return trackGainers, nil
}

func extractPriceData(ctx context.Context, client graphql.Client, trackGainers *[]shared.Gainers, datetime, long int, botName string) ([]SMA, error) {

	var allPriceData []SMA
	// Initialize the rolling totals map
	rollingTotals := make(map[string]RollingTotal)

	for timeFrameCount := 0; timeFrameCount <= long; timeFrameCount++ {
		var timeFramePrices SMA
		targetDatetime := datetime - (timeFrameCount * 300)
		log.Info().Str("Bot", botName).Int("Supplied datetime", datetime).Int("time Frame Count", timeFrameCount).Int("target Datetime", targetDatetime).Msg("Gathering price data for")
		historicPricesList, err := GetPriceData(ctx, client, targetDatetime, strconv.Itoa(timeFrameCount)+botName)
		if err != nil {
			log.Error().Msgf("CompareSMAs - historicPricesList!")
			return nil, err
		}

		// Initialize the priceDataArray here, outside the timeFrameCount loop
		priceDataArray := make([]model.Pair, len(*trackGainers))
		var Price float64

		for i := 0; i < len(historicPricesList); i++ {
			for j, pair := range *trackGainers {
				coin := historicPricesList[i].Symbol
				if coin == pair.Symbol {
					log.Info().Str("coin", coin).Str("Symbol", pair.Symbol).Msg("start")

					// Convert coin.Price to a float64
					Price, err = strconv.ParseFloat(historicPricesList[i].Price, 64)
					if err != nil {
						log.Error().Msgf("CompareSMAs - failed to parse coin price as float64: %v", err)
						return nil, err
					}

					// Increment the rolling total for this coin
					if rt, found := rollingTotals[pair.Symbol]; found {
						rollingTotals[pair.Symbol] = RollingTotal{
							Symbol:       pair.Symbol,
							RollingTotal: rt.RollingTotal + Price,
						}
					} else {
						rollingTotals[pair.Symbol] = RollingTotal{
							Symbol:       pair.Symbol,
							RollingTotal: Price,
						}
					}
					// Update the priceDataArray
					priceDataArray[j] = model.Pair{Symbol: pair.Symbol, Price: historicPricesList[i].Price}
					timeFramePrices = SMA{TimeFrame: timeFrameCount, PriceDataArray: priceDataArray}

				}
			}
		}

		allPriceData = append(allPriceData, timeFramePrices)

	}

	return allPriceData, nil
}

func ProcessAllPriceData(allPriceData []SMA, short, long int, botName string) (map[string]float64, map[string]float64, map[string]float64, error) {

	fmt.Println("INSIDE PROCESS ALL PRICE DATA")
	fmt.Println(allPriceData)
	// Reverse the array
	rollingTotals := make(map[string]float64)
	currentPrice, shortAve, longAve := make(map[string]float64), make(map[string]float64), make(map[string]float64)
	log.Info().Str("Bot", botName).Int("Short", short).Int("Long", long).Int("prices", len(allPriceData)).Msg("Calculating Simple moving averages ... ")
	for _, timeFramePrices := range allPriceData {
		for _, priceDataArray := range timeFramePrices.PriceDataArray {
			symbol := priceDataArray.Symbol
			price, err := strconv.ParseFloat(priceDataArray.Price, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			if _, ok := rollingTotals[symbol]; !ok {
				rollingTotals[symbol] = 0
			}
			rollingTotals[symbol] += price

			//if timeFramePrices.TimeFrame == 0 {
			currentPrice[symbol] = price // Update the currentPrice map
			log.Info().Str("Bot", botName).Str("Symbol", symbol).Float64("Price", currentPrice[symbol]).Msg("Current Price")
			//}

			if timeFramePrices.TimeFrame == short-1 {
				shortAve[symbol] = rollingTotals[symbol] / float64(short)
				shortAverage := (shortAve[symbol])
				log.Info().Int("Short", short-1).Int("timefream if", timeFramePrices.TimeFrame).Str("Bot", botName).Int("Minutes", short*5).Str("Symbol", symbol).Float64("Average", shortAverage).Msg("Short period Simple Moving Ave")
			}

			if timeFramePrices.TimeFrame == long-1 {
				longAve[symbol] = rollingTotals[symbol] / float64(long)
				longAverage := longAve[symbol]
				log.Info().Str("Bot", botName).Int("Minutes", long*5).Str("Symbol", symbol).Float64("Average", longAverage).Msg("Long period Simple Moving Ave")
			}
			//fmt.Println("A")
		}
		//fmt.Println("B")
	}
	//fmt.Println("C")
	return currentPrice, shortAve, longAve, nil
}
