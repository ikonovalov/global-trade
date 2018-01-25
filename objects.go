/*
 * MIT License
 *
 * Copyright (c) 2018 Igor Konovalov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"encoding/json"
	"fmt"
)

// ============ PUBLIC API OBJECTS =================
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
	Offers map[string]Offers
}

type Offers struct {
	Asks []Offer `json:"asks"`
	Bids []Offer `json:"bids"`
}

type Offer struct {
	Price    float64
	Quantity float64
}

func (n *Offer) UnmarshalJSON(buf []byte) error {
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

/* trades
{
	"ltc_btc":[
		{
			"type":"ask",
			"price":104.2,
			"amount":0.101,
			"tid":41234426,
			"timestamp":1418654531
		},
		{
			"type":"bid",
			"price":103.53,
			"amount":1.51414,
			"tid":41234422,
			"timestamp":1418654530
		},
		...
	]
	...
}
 */
type TradesResponse struct {
	Trades map[string][]Trade
}

type Trade struct {
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
	Tid       uint64  `json:"tid"`
	Timestamp int64   `json:"timestamp"`
}

// ============ TRADE API OBJECTS =====================
type GetInfoResponse struct {
	Success uint8         `json:"success"`
	Data    GetInfoReturn `json:"return"`
	Error   string        `json:"error"`
}

type GetInfoReturn struct {
	Rights             map[string]uint8   `json:"rights"`
	Funds              map[string]float64 `json:"funds"`
	FundsIncludeOrders map[string]float64 `json:"funds_incl_orders"`
	TransactionCount   int                `json:"transaction_count"`
	OpenOrders         int                `json:"open_orders"`
	ServerTime         int64              `json:"server_time"`
}

type ActiveOrdersResponse struct {
	Success uint8                  `json:"success"`
	Orders  map[string]ActiveOrder `json:"return"`
	Error   string                 `json:"error"`
}

type ActiveOrder struct {
	Pair    string  `json:"pair"`
	Type    string  `json:"type"`
	Amount  float64 `json:"amount"`
	Rate    float64 `json:"rate"`
	Created string  `json:"timestamp_created"`
	Status  uint8   `json:"status"`
}

type OrderInfoResponse struct {
	Success uint8                `json:"success"`
	Orders  map[string]OrderInfo `json:"return"`
	Error   string               `json:"error"`
}

type OrderInfo struct {
	Pair        string  `json:"pair"`
	Type        string  `json:"type"`
	StartAmount float64 `json:"start_amount"`
	Amount      float64 `json:"amount"`
	Rate        float64 `json:"rate"`
	Created     string  `json:"timestamp_created"`
	Status      uint8   `json:"status"`
}

type TradeHistoryResponse struct {
	Success uint8                    `json:"success"`
	Error   string                   `json:"error"`
	Orders  map[string]HistoricOrder `json:"return"`
}

type HistoricOrder struct {
	Pair        string  `json:"pair"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Rate        float64 `json:"rate"`
	OrderId     uint64  `json:"order_id"`
	IsYourOrder uint8   `json:"is_your_order"`
	Timestamp   uint64  `json:"timestamp"`
}

type TradeResponse struct {
	Success uint8       `json:"success"`
	Error   string      `json:"error"`
	Result  TradeResult `json:"return"`
}

type TradeResult struct {
	Received float64            `json:"received"`
	Remains  float64            `json:"remains"`
	OrderId  string             `json:"order_id"`
	Funds    map[string]float64 `json:"funds"`
}

type CancelOrderRespose struct {
	Success uint8        `json:"success"`
	Error   string       `json:"error"`
	Result  CancelResult `json:"return"`
}

type CancelResult struct {
	OrderId string             `json:"order_id"`
	Funds   map[string]float64 `json:"funds"`
}
