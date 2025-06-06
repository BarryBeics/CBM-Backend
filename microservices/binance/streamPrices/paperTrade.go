package binanace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"cryptobotmanager.com/cbm-backend/microservices/reports/functions"
	"cryptobotmanager.com/cbm-backend/shared"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Stream payload format from Binance
type BinanceTrade struct {
	Price     string `json:"p"`
	TradeTime int64  `json:"T"`
	Symbol    string `json:"s"`
}

// TradeValues holds calculated target thresholds
type TradeValues struct {
	TimedOut   int64
	TakeProfit float64
	StopLoss   float64
}

func ListenAndPaperTrade(ctx context.Context, client graphql.Client, symbol string, details model.StrategyInput) {
	lowerSymbol := strings.ToLower(symbol)
	wsURL := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", lowerSymbol)

	botName := details.BotInstanceName
	accountBalance := details.AccountBalance
	feesBalance := details.FeesTotal

	// Connect
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("WebSocket dial failed")
		return
	}
	defer conn.Close()

	log.Info().Str("symbol", symbol).Msg("WebSocket connection opened")
	startTime := time.Now().UnixMilli()

	// LIVE DATA - Get latest price
	openingPrice, err := getLatestPrice(symbol)
	if err != nil {
		fmt.Println("Error getting latest price:", err)
		return
	}

	exitValues := calculateExitValues(openingPrice, details, int(startTime))
	log.Info().
		Float64("latest price", openingPrice).
		Int("The trade will time out at:", exitValues.TimedOut).
		Float64("Take profit set at:", exitValues.TakeProfit).
		Float64("Stop Loss set at", exitValues.StopLoss).
		Msg("Exits")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error().
				Err(err).
				Str("symbol", symbol).
				Msg("WebSocket read failed")
			return
		}

		var tradeUpdate BinanceTrade
		if err := json.Unmarshal(message, &tradeUpdate); err != nil {
			log.Error().
				Err(err).
				Str("raw", string(message)).
				Msg("Failed to decode trade message")
			continue
		}

		currentPrice, err := strconv.ParseFloat(tradeUpdate.Price, 64)
		if err != nil {
			log.Warn().
				Str("price", tradeUpdate.Price).
				Msg("Price parse error")
			continue
		}

		var volume float64

		// Evaluate exit conditions
		streamTime := int(tradeUpdate.TradeTime)
		elapsedTime := streamTime - int(startTime)

		outCome := graph.UpdateCountersInput{}
		var TIMEOUTGainCounter, TIMEOUTLossCounter bool

		switch {
		case currentPrice >= exitValues.TakeProfit:
			change := shared.PercentageChange(openingPrice, currentPrice)
			updatedBalance, fees, netOutcome := CalculateUpdatedBalance(accountBalance, change, 0.06)
			functions.TradeOutcomeReport(client, streamTime, elapsedTime, botName, change, updatedBalance, volume, fees, symbol, "WIN")

			outCome = graph.UpdateCountersInput{
				BotInstanceName:    botName,
				WINCounter:         true,
				LOSSCounter:        false,
				TIMEOUTGainCounter: false,
				TIMEOUTLossCounter: false,
				NetGainCounter:     netOutcome,
				NetLossCounter:     !netOutcome,
				AccountBalance:     updatedBalance,
				FeesTotal:          *feesBalance + fees}
			graph.UpdateCounters(ctx, client, outCome)
			log.Info().
				Float64("% Change", change).Int64("Start Time", startTime).
				Int("Stream Time", streamTime).Int("elapsed Time", elapsedTime).
				Float64("Open price", openingPrice).
				Float64("Take profit", exitValues.TakeProfit).
				Str("symbol", symbol).
				Msg("We have WON on this trade")

				// Close the WebSocket connection
			conn.Close()
			return

		case currentPrice <= exitValues.StopLoss:
			change := shared.PercentageChange(openingPrice, currentPrice)
			updatedBalance, fees, _ := CalculateUpdatedBalance(accountBalance, change, 0.06)
			outCome = graph.UpdateCountersInput{
				BotInstanceName:    botName,
				WINCounter:         false,
				LOSSCounter:        true,
				TIMEOUTGainCounter: false,
				TIMEOUTLossCounter: false,
				AccountBalance:     updatedBalance,
				FeesTotal:          *feesBalance + fees}
			graph.UpdateCounters(ctx, client, outCome)
			log.Info().
				Str("symbol", symbol).
				Float64("price", currentPrice).
				Float64("% Change", change).
				Int("Elapsed(ms)", elapsedTime).
				Msg("STOP LOSS hit")

				// Close the WebSocket connection
			conn.Close()
			return

		case streamTime >= exitValues.TimedOut:
			change := shared.PercentageChange(openingPrice, currentPrice)
			updatedBalance, fees, netOutcome := CalculateUpdatedBalance(accountBalance, change, 0.06)
			functions.TradeOutcomeReport(client, streamTime, elapsedTime, botName, change, updatedBalance, volume, fees, symbol, "TIMED OUT")

			// Assuming change is the percentage change
			if change > 0 {
				TIMEOUTGainCounter = true
			} else if change < 0 {
				TIMEOUTLossCounter = true
			}

			winCounter := false
			lossCounter := false
			feesTotal := *feesBalance + fees
			outCome := graph.UpdateCountersInput{
				BotInstanceName:    botName,
				WINCounter:         winCounter,
				LOSSCounter:        lossCounter,
				TIMEOUTGainCounter: TIMEOUTGainCounter,
				TIMEOUTLossCounter: TIMEOUTLossCounter,
				NetGainCounter:     netOutcome,
				NetLossCounter:     !netOutcome,
				AccountBalance:     updatedBalance,
				FeesTotal:          feesTotal,
			}

			// Call the UpdateCounters function with the updated outCome
			graph.UpdateCounters(ctx, client, outCome)
			log.Info().
				Float64("% Change", change).
				Int64("Start Time", startTime).
				Int("Stream Time", streamTime).
				Int("elapsed Time", elapsedTime).
				Float64("Open price", openingPrice).
				Int("Timed out at", exitValues.TimedOut).
				Str("symbol", symbol).
				Msg("We have been TIMED on this trade")

			// Close the WebSocket connection
			conn.Close()

			return
		}
	}
}

func getLatestPrice(symbol string) (float64, error) {
	key := "https://api.binance.com/api/v3/ticker/price?symbol=" + strings.ToUpper(symbol)
	resp, err := http.Get(key)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	price, ok := data["price"].(string)
	if !ok {
		return 0, fmt.Errorf("price not found in response")
	}

	binanceStreamPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, err
	}

	return binanceStreamPrice, nil
}

func calculateExitValues(price float64, details model.StrategyInput, startTime int) shared.TradeValues {
	tradeDuration := details.TradeDuration * 60000
	timedOut := (startTime + int(tradeDuration))

	// Assuming details.TakeProfitPercentage and details.StopLossPercentage are in decimal form (e.g., 1.5% as 0.015)
	takeProfitFactor := 1 + (details.TakeProfitPercentage / 100)
	stopLossFactor := 1 - (details.StopLossPercentage / 100)

	takeProfit := price * takeProfitFactor
	stopLoss := price * stopLossFactor

	log.Debug().
		Int("Seconds trade will time out in:", timedOut).
		Float64("Trade duration:", float64(tradeDuration)).
		Float64("Start Time:", float64(startTime)).
		Float64("Take profit set at:", takeProfit).
		Float64("Stop Loss set at", stopLoss).
		Msg("Exits")
	return shared.TradeValues{
		TimedOut:   timedOut,
		TakeProfit: takeProfit,
		StopLoss:   stopLoss,
	}
}

// CalculateUpdatedBalance calculates the updated balance after applying the percentage change and subtracting the fees.
func CalculateUpdatedBalance(balance, change, feePercentage float64) (float64, float64, bool) {

	log.Debug().
		Float64("Balance", balance).
		Msg("Opening Balance")
	// Calculate the fees for converting BNB to the other coin
	buyFee := balance * feePercentage / 100

	// Subtract the fees for converting BNB to the other coin
	balance -= buyFee
	log.Debug().
		Float64("Buy Fee", buyFee).
		Float64("updated Balance", balance).
		Msg("After Buy")

	// Calculate the balance after applying the percentage change
	updatedBalance := balance + (balance * change / 100)
	log.Debug().
		Float64("Change %", change).
		Float64("New Value", updatedBalance).
		Msg("after trade")

	// Calculate the fees for converting the other coin back to BNB
	sellFee := updatedBalance * feePercentage / 100

	// Subtract the fees for converting the other coin back to BNB
	updatedBalance -= sellFee
	log.Debug().
		Float64("Sell Fee", sellFee).
		Float64("updated Balance", updatedBalance).
		Msg("After Sell")

	totalFees := buyFee + sellFee
	log.Debug().Float64("fees", totalFees).Msg("Total")

	gainLoss := updatedBalance - balance

	netGain := false

	if gainLoss > 0 {
		netGain = true
		log.Debug().Bool("Net Gain", netGain).Msg("Outcome")
	} else {
		log.Debug().Bool("Net Gain", netGain).Msg("Outcome")
	}

	return updatedBalance, totalFees, netGain
}
