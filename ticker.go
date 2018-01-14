package main

type TickerInfoResponse struct {
	Tickers map[string]Ticker
}

type Ticker struct {
	High float64 `json:"high"`
	Low  float64 `json:"low"`
	Avg  float64 `json:"avg"`
	Vol  float64 `json:"vol"`
	VolCur  float64 `json:"vol_cur"`
	Buy  float64 `json:"buy"`
	Sell  float64 `json:"sell"`
	Updated  int64 `json:"updated"`
}