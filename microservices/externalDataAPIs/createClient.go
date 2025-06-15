package binance

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog/log"
)

// NewBinanceClient creates and returns a new Binance API client using environment variables.
func NewBinanceClient() *binance.Client {
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")
	client := binance.NewClient(apiKey, secretKey)
	if client == nil {
		slog.Error("Failed to create Binance client", "error", fmt.Errorf("failed to create Binance client"))
		return nil
	}
	log.Info().Msg("Created connection to Binance")
	return client
}
