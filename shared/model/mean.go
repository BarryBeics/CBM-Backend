package model

type MeanStat struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}

type MeanStatInput struct {
	Avg   float64 `json:"Avg"`
	Count int     `json:"Count"`
}
