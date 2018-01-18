package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
)

var (
	command   = kingpin.Arg("command", "Command").Required().String()
	ticker   = kingpin.Arg("ticker", "Ticker").String()
)

func tickers24(client *http.Client, pairs string, ch chan TickerInfoResponse) {
	url := ApiBase + ApiVersion + "/ticker/" + pairs
	response := query(client, url)

	var tickerResponse TickerInfoResponse
	pTicker := &tickerResponse.Tickers

	if err := unmarshal(response, pTicker); err != nil {
		panic(err)
	}
	ch <- tickerResponse
}

func info(client *http.Client, ch chan InfoResponse) {
	url := ApiBase + ApiVersion + "/info"
	response := query(client, url)
	var infoResponse InfoResponse
	if err := unmarshal(response, &infoResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- infoResponse
}

func depth(client *http.Client, pairs string, ch chan DepthResponse)  {
	depthLimited(client, pairs, 150, ch)
}

func depthLimited(client *http.Client, pairs string, limit uint8, ch chan DepthResponse)  {
	url := ApiBase + ApiVersion + "/depth/" + pairs + "?limit=" + string(limit)
	response := query(client, url)
	var depthResponse DepthResponse
	depthResponse.Raw = string(response)
	if err := unmarshal(response, &depthResponse.Orders); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- depthResponse
}

func unmarshal(data [] byte, obj interface{}) error {
	error := json.Unmarshal(data, obj)
	if error != nil {
		log.Fatal("Unmarshaling failed\n" + string(data))
	}
	return error
}

func query(client *http.Client, url string) ([]byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Errorf("error. HTTP %d", resp.StatusCode)
		panic(errors.New("Something goes wrong. We got HTTP" + string(resp.StatusCode)))
	}
	response, _ := ioutil.ReadAll(resp.Body)
	return response
}

type Yobit struct {
	client *http.Client
}