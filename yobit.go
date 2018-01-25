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
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"strconv"
	"net/url"
	"bytes"
	"github.com/ikonovalov/go-cloudflare-scraper"
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
	ApiTrade   = "https://yobit.net/tapi/"
)

func NewYobit() (*Yobit) {
	cloudflare, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		panic(err)
	}

	keys, err := loadApiKeys()
	if err != nil {
		panic(err)
	}

	yobit := Yobit{client: &http.Client{Transport: cloudflare, Jar: cloudflare.Cookies}, apiKeys: &keys}
	yobit.PassCloudflare()
	return &yobit
}

func (y *Yobit) PassCloudflare() (*Yobit) {
	channel := make(chan TickerInfoResponse)
	go y.Tickers24("btc_usd", channel)
	<-channel
	return y
}

func (y *Yobit) Tickers24(pairs string, ch chan<- TickerInfoResponse) {
	ticker24Url := ApiBase + ApiVersion + "/ticker/" + pairs
	response := y.callPublic(ticker24Url)

	var tickerResponse TickerInfoResponse
	pTicker := &tickerResponse.Tickers

	if err := unmarshal(response, pTicker); err != nil {
		panic(err)
	}
	ch <- tickerResponse
}

func (y *Yobit) Info(ch chan<- InfoResponse) {
	infoUrl := ApiBase + ApiVersion + "/info"
	response := y.callPublic(infoUrl)
	var infoResponse InfoResponse
	if err := unmarshal(response, &infoResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- infoResponse
}

func (y *Yobit) Depth(pairs string, ch chan<- DepthResponse) {
	y.DepthLimited(pairs, 150, ch)
}

func (y *Yobit) DepthLimited(pairs string, limit int, ch chan<- DepthResponse) {
	limitedDepthUrl := fmt.Sprintf("%s/depth/%s?limit=%d", ApiBase+ApiVersion, pairs, limit)
	response := y.callPublic(limitedDepthUrl)
	var depthResponse DepthResponse
	if err := unmarshal(response, &depthResponse.Offers); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- depthResponse
}

func (y *Yobit) TradesLimited(pairs string, limit int, ch chan<- TradesResponse) {
	tradesLimitedUrl := fmt.Sprintf("%s/trades/%s?limit=%d", ApiBase+ApiVersion, pairs, limit)
	response := y.callPublic(tradesLimitedUrl)
	var tradesResponse TradesResponse
	if err := unmarshal(response, &tradesResponse.Trades); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- tradesResponse
}

// TRADE API =================================================================================

func (y *Yobit) GetInfo(ch chan<- GetInfoResponse) {
	response := y.callPrivate("getInfo")
	var getInfoResp GetInfoResponse
	if err := unmarshal(response, &getInfoResp); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if getInfoResp.Success == 0 {
		panic(errors.New(getInfoResp.Error))
	}
	ch <- getInfoResp
}

func (y *Yobit) ActiveOrders(pair string, ch chan<- ActiveOrdersResponse) {
	response := y.callPrivate("ActiveOrders", CallArg{"pair", pair})
	var activeOrders ActiveOrdersResponse
	if err := unmarshal(response, &activeOrders); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if activeOrders.Success == 0 {
		panic(errors.New(activeOrders.Error))
	}
	ch <- activeOrders
}

func (y *Yobit) OrderInfo(orderId string, ch chan<- OrderInfoResponse) {
	response := y.callPrivate("OrderInfo", CallArg{"order_id", orderId})
	var orderInfo OrderInfoResponse
	if err := unmarshal(response, &orderInfo); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if orderInfo.Success == 0 {
		panic(errors.New(orderInfo.Error))
	}
	ch <- orderInfo
}

func (y *Yobit) Trade(pair string, tradeType string, rate float64, amount float64, ch chan TradeResponse) {
	response := y.callPrivate("Trade",
		CallArg{"pair", pair},
		CallArg{"type", tradeType},
		CallArg{"rate", strconv.FormatFloat(rate, 'f', 8, 64)},
		CallArg{"amount", strconv.FormatFloat(amount, 'f', 8, 64)},
	)
	var tradeResponse TradeResponse
	if err := unmarshal(response, &tradeResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if tradeResponse.Success == 0 {
		panic(errors.New(tradeResponse.Error))
	}
	ch <- tradeResponse
}

func (y *Yobit) CancelOrder(orderId string , ch chan CancelOrderRespose) {
	response := y.callPrivate("CancelOrder", CallArg{"order_id", orderId})
	var cancelResponse CancelOrderRespose
	if err := unmarshal(response, &cancelResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if cancelResponse.Success == 0 {
		panic(errors.New(cancelResponse.Error))
	}
	ch <- cancelResponse
}

func (y *Yobit) TradeHistory(pair string, ch chan<- TradeHistoryResponse) {
	response := y.callPrivate("TradeHistory",
		CallArg{"pair", pair},
		CallArg{"count", "1000"},
	)
	var tradeHistory TradeHistoryResponse
	if err := unmarshal(response, &tradeHistory); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if tradeHistory.Success == 0 {
		panic(errors.New(tradeHistory.Error))
	}
	ch <- tradeHistory
}

func unmarshal(data [] byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unmarshaling failed\n%s\n%s", string(data), err.Error()))
	}
	return err
}

func (y *Yobit) query(req *http.Request) ([]byte) {
	resp, err := y.client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("error. HTTP %d", resp.StatusCode)
		panic(errors.New(fmt.Sprintf("\n%s\nSomething goes wrong. We got HTTP%d", req.URL.String(), resp.StatusCode)))
	}
	response, _ := ioutil.ReadAll(resp.Body)
	return response
}

func (y *Yobit) callPublic(url string) ([]byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic(err)
	}
	return y.query(req)
}

type CallArg struct {
	name, value string
}

func (y *Yobit) callPrivate(method string, args ...CallArg) ([]byte) {
	nonce := getNonce()
	form := url.Values{
		"method": {method},
		"nonce":  {strconv.FormatUint(nonce, 10)},
	}
	for _, arg := range args {
		form.Add(arg.name, arg.value)
	}
	encode := form.Encode()
	signature := signHmacSha512([]byte(y.apiKeys.Secret), []byte(encode))
	body := bytes.NewBufferString(encode)
	req, err := http.NewRequest("POST", ApiTrade, body)

	if err != nil {
		log.Fatal("callPrivate: ", err)
		panic(err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", y.apiKeys.Key)
	req.Header.Add("Sign", signature)

	query := y.query(req)
	incrementNonce(&nonce)
	return query
}

type Yobit struct {
	client  *http.Client
	apiKeys *ApiKeys
}
