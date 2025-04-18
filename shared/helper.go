package shared

import (
	"math"
	"time"
)

// RoundMinuteToFiveMinuteInterval function takes in a currentMinute integer as an
// argument, rounds it to the nearest five-minute interval and returns the rounded
// minute as an integer. It also logs debug messages using zerolog package.
// RoundTimeToFiveMinuteInterval takes in a time.Time value, rounds it to the nearest five-minute interval,
// and returns the rounded time.
func RoundTimeToFiveMinuteInterval(epochTime int64) int {
	// Convert epoch time to time.Time
	t := time.Unix(epochTime, 0)

	// Round to the nearest 5 minutes
	rounded := t.Round(5 * time.Minute)

	// Convert the rounded time back to epoch time
	int64Value := rounded.Unix()

	// Log rounded time
	//log.Trace().Int64("original_epoch_time", epochTime).Int64("rounded_epoch_time", int64Value).Msg("rounded time to five-minute interval")

	var roundedEpochTime int = int(int64Value)

	return roundedEpochTime
}

// RoundFloatToDecimal rounds the given floating-point number 'x'
// to the specified number of decimal places.
func RoundFloatToDecimal(x, decimal float64) float64 {
	if decimal < 1 || decimal > 4 {
		return x // Do not round if decimal is not 1, 2, or 3
	}
	unit := CalculateMultiplier(decimal)
	return math.Round(x*unit) / unit
}

// CalculateMultiplier returns the multiplier corresponding to the given number.
// For example, when 'num' is 1, it returns 10; when 'num' is 2, it returns 100; and so on.
func CalculateMultiplier(num float64) float64 {
	if num == 1 {
		return 10
	} else if num == 2 {
		return 100
	} else if num == 3 {
		return 1000
	} else if num == 4 {
		return 10000
	} else {
		return 0
	}
}

// FindUniqueStrings returns the unique strings present in slice1 but not in slice2.
func FindUniqueStrings(slice1, slice2 []string) []string {
	unique := make(map[string]struct{})

	// Add items from slice1 to unique map
	for _, item := range slice1 {
		unique[item] = struct{}{}
	}

	// Remove items from slice2 that are already in the unique map
	for _, item := range slice2 {
		delete(unique, item)
	}

	// Convert the keys of the unique map to a slice
	uniqueItems := make([]string, 0, len(unique))
	for item := range unique {
		uniqueItems = append(uniqueItems, item)
	}

	return uniqueItems
}

// Percentage calculates the percentage of inputTwo relative to inputOne.
// It returns the calculated percentage rounded to two decimal places.
func Percentage(inputOne, inputTwo float64) (result float64) {
	share := 100 / inputOne
	result = share * inputTwo
	result = result - 100

	return RoundFloatToDecimal(result, 2)
}

func Round(x, decimal float64) float64 {
	unit := calculateValue(decimal)
	return math.Round(x*unit) / unit
}

// not sure what this is for - i will probably discover it when I rebuild the application
func calculateValue(num float64) float64 {
	if num == 1 {
		return 10
	} else if num == 2 {
		return 100
	} else if num == 3 {
		return 1000
	} else {
		return 0
	}
}
