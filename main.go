package main

import (
	"encoding/json"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"os"
	"strings"
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
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

func main() {

	argsWithProg := os.Args
	var pair string
	if len(argsWithProg) == 2 {
		pair = argsWithProg[1]
	} else {
		pair = "btc_usd-ltc_usd-xem_usd-eth_usd-xem_eth"
	}

	client := &http.Client{}
	ch := make(chan TickerInfoResponse)
	go tickers24(client, pair, ch)
	tickerResponse := <-ch

	// print result
	for ticker, v := range tickerResponse.Tickers {
		fmt.Printf(
			"%-9s High [%.8f] Avg [%.8f] Low[%.8f] Volume[%.8f] Current Volume[%.8f] Buy[%.8f] Sell[%.8f] Last[%.8f]\n",
			Bold(strings.ToUpper(ticker)), v.High, Green(v.Avg), v.Low, v.Vol, v.VolCur, v.Buy, v.Sell, v.Last)
	}

	ch2 := make(chan DepthResponse)
	go depth(client, pair, ch2)
	depthResponse := <- ch2
	orders := depthResponse.Orders[pair]
	for idx, ask := range orders.Asks {
		fmt.Printf("#%d %.8f <- %.8f\n" , idx, ask.Price, ask.Quantity)
	}

}


