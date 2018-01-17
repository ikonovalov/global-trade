package main

import (
	"encoding/json"
	"fmt"
)

/** ticker
{
	"ltc_btc":{
		"high":105.41,
		"low":104.67,
		"avg":105.04,
		"vol":43398.22251455,
		"vol_cur":4546.26962359,
		"last":105.11,
		"buy":104.2,
		"sell":105.11,
		"updated":1418654531
	}
	...
}
 */
type TickerInfoResponse struct {
	Tickers map[string]Ticker
}

type Ticker struct {
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Avg     float64 `json:"avg"`
	Vol     float64 `json:"vol"`
	VolCur  float64 `json:"vol_cur"`
	Buy     float64 `json:"buy"`
	Sell    float64 `json:"sell"`
	Last    float64 `json:"last"`
	Updated int64   `json:"updated"`
}

/** info
	{
	"server_time":1418654531,
	"pairs":{
		"ltc_btc":{
			"decimal_places":8,
			"min_price":0.00000001,
			"max_price":10000,
			"min_amount":0.0001,
			"hidden":0,
			"fee":0.2
		}
		...
	}
}
 */
type InfoResponse struct {
	ServerTime int64               `json:"server_time"`
	Pairs      map[string]PairInfo `json:"pairs"`
}

type PairInfo struct {
	DecimalPlace uint16  `json:"decimal_places"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	MinAmount    float64 `json:"min_amount"`
	Hidden       uint8   `json:"hidden"`
	Fee          float64 `json:"fee"`
}

/** depth
{
	"ltc_btc":{
		"asks":[
			[104.67,0.01],
			[104.75,11],
			[104.80,0.523],
			...
		],
		"bids":[
			[104.3,5.368783],
			[104.212,2.57357],
			[103.62,0.43663336],
			[103.61,0.7255672],
			...
		]
	}
	...
}
 */
type DepthResponse struct {
	Raw    string
	Orders map[string]Orders
}

type Orders struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
}

type Order struct {
	Price float64
	Quantity float64
}

func (n *Order) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Price, &n.Quantity}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Order: %d != %d", g, e)
	}
	return nil
}
