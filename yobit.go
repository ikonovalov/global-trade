package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"crypto/hmac"
	"crypto/sha512"
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
)

func (y *Yobit) Tickers24(pairs string, ch chan TickerInfoResponse) {
	url := ApiBase + ApiVersion + "/ticker/" + pairs
	response := y.query(url)

	var tickerResponse TickerInfoResponse
	pTicker := &tickerResponse.Tickers

	if err := y.unmarshal(response, pTicker); err != nil {
		panic(err)
	}
	ch <- tickerResponse
}

func (y *Yobit) Info(ch chan InfoResponse) {
	url := ApiBase + ApiVersion + "/info"
	response := y.query(url)
	var infoResponse InfoResponse
	if err := y.unmarshal(response, &infoResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- infoResponse
}

func (y *Yobit) Depth(pairs string, ch chan DepthResponse) {
	y.DepthLimited(pairs, 150, ch)
}

func (y *Yobit) DepthLimited(pairs string, limit int, ch chan DepthResponse) {
	url := fmt.Sprintf("%s/depth/%s?limit=%d", ApiBase + ApiVersion, pairs, limit)
	response := y.query(url)
	var depthResponse DepthResponse
	if err := y.unmarshal(response, &depthResponse.Orders); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- depthResponse
}

func (y *Yobit) TradesLimited(pairs string, limit int, ch chan TradesResponse) {
	url := fmt.Sprintf("%s/trades/%s?limit=%d", ApiBase + ApiVersion, pairs, limit)
	response := y.query(url)
	var tradesResponse TradesResponse
	if err := y.unmarshal(response, &tradesResponse.Trades); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- tradesResponse
}

func (y *Yobit) unmarshal(data [] byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		log.Fatal("Unmarshaling failed\n" + string(data))
	}
	return err
}

func (y *Yobit) query(url string) ([]byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic(err)
	}

	resp, err := y.client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Errorf("error. HTTP %d", resp.StatusCode)
		panic(errors.New(fmt.Sprintf("\n%s\nSomething goes wrong. We got HTTP%d", url, resp.StatusCode)))
	}
	response, _ := ioutil.ReadAll(resp.Body)
	return response
}

func (y *Yobit) sign(key []byte, message []byte) ([]byte) {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	digest := mac.Sum(nil)
	return digest
}

type Yobit struct {
	client *http.Client
}
