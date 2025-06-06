package shared

type Gainers struct {
	IncrementPriceGain float64
	SMAPriceGain       float64
	Symbol             string
	ATR                float64
	WeightedScore      float64
}

// Define a struct to hold the trade values
type TradeValues struct {
	TimedOut   int
	TakeProfit float64
	StopLoss   float64
}
