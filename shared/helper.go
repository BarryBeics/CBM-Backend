package shared

import "time"

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
