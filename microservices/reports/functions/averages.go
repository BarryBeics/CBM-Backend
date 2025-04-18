package functions

import (
	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/go-gota/gota/dataframe"
	"github.com/rs/zerolog/log"
)

// AverageGain function takes in a slice of Gainers and an integer value topN as
// arguments. It converts the slice into a dataframe, sorts it in descending order
// by IncrementPriceGain, and then selects the top N coins using Subset() method.
// It then calculates the mean of IncrementPriceGain for the selected coins and
// logs relevant information using zerolog package. The function returns the
// average gain as a float64 value.
func AverageGain(coinsWithMomentum *[]shared.Gainers, topN int) float64 {
	newOne := *coinsWithMomentum

	qty := len(newOne)
	coinsDf := dataframe.LoadStructs(newOne)
	//fmt.Print("COINS WITH MOMENTUM", coinsDf)

	sorted := coinsDf.Arrange(dataframe.RevSort("IncrementPriceGain"))
	//fmt.Print("SORTED", sorted)

	// Ensure topN does not exceed the number of coins
	if topN > qty {
		topN = qty
	}

	// Check if the DataFrame is empty
	if qty == 0 {
		log.Warn().Msg("DataFrame is empty")
		return 0.0
	}

	// Check if topN is valid
	if topN <= 0 {
		log.Warn().Msg("Invalid topN value")
		return 0.0
	}

	// Create a subset only if there are elements in the DataFrame
	subsetIndices := make([]int, topN)
	for i := 0; i < topN; i++ {
		subsetIndices[i] = i
	}

	sub := sorted.Subset(subsetIndices)
	//fmt.Print("SUBSET", sub)

	// Log top N coins with momentum
	log.Trace().Int("top_n", topN).Msg("top coins with momentum")

	averageDF := sub.Col("IncrementPriceGain")
	col := coinsDf.Col("Symbol")
	tickers := col.Records()

	average := averageDF.Mean()

	// Log average gain for top N coins
	log.Debug().Int("coins_qty", qty).Float64("avg_gain", average).Interface("tickers", tickers).Msg("average gain for top coins")

	return shared.Round(average, 2)
}
