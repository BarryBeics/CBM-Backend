package functions

// import (
// 	"context"
// 	"strconv"

// 	"cryptobotmanager.com/cbm-backend/shared/binance"
// )

// func FetchLiquidityData() (map[string]float64, error) {

// 	client := binance.NewBinanceClient()

// 	stats, err := client.NewListPriceChangeStatsService().Do(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	liquidity := make(map[string]float64)
// 	for _, s := range stats {
// 		vol, err := strconv.ParseFloat(s.QuoteVolume, 64)
// 		if err != nil {
// 			continue
// 		}
// 		liquidity[s.Symbol] = vol
// 	}

// 	return liquidity, nil
// }
