package model

type Mean struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}

type MeanInput struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}
