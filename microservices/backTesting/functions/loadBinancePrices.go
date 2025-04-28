package functions

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
)

func ExtractTimestampFromFilename(filename string) int {
	base := filepath.Base(filename)
	parts := strings.TrimSuffix(strings.TrimPrefix(base, "binance_prices_"), ".json")

	// Parse date like 2025-03-30
	t, err := time.Parse("2006-01-02", parts)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to parse date from filename")
		return int(time.Now().Unix())
	}
	// Add a fake time (e.g. noon)
	t = t.Add(12 * time.Hour)
	return int(t.Unix())
}

func LoadPriceSnapshotsFromFile(path string) ([]model.NewHistoricPriceInput, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var snapshots []model.NewHistoricPriceInput
	err = json.Unmarshal(bytes, &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
