package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FNGResponse struct {
	Data []struct {
		Value               string `json:"value"`
		ValueClassification string `json:"value_classification"`
		Timestamp           string `json:"timestamp"`
	} `json:"data"`
}

func main() {
	resp, err := http.Get("https://api.alternative.me/fng/?limit=1")
	if err != nil {
		log.Fatal("Failed to fetch FNG index:", err)
	}
	defer resp.Body.Close()

	var result FNGResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to decode response:", err)
	}

	if len(result.Data) == 0 {
		log.Fatal("No data returned from FNG API")
	}

	latest := result.Data[0]
	fmt.Printf("Fear & Greed Index: %s (%s) at %s\n", latest.Value, latest.ValueClassification, latest.Timestamp)

	// TODO: persist to DB, write to file, or include in your trade outcome report
}
