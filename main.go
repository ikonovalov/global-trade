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
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
)

func tickers(pairs string) string {
	return ApiBase + ApiVersion + "/ticker/" + pairs
}

func unmarshal(data [] byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
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
		panic(errors.New("Status is not 200. It's " + string(resp.StatusCode)))
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
		pair = "xem_usd-eth_usd-xem_eth"
	}

	url := tickers(pair)
	fmt.Println(Bold(url + "\n"))

	client := &http.Client{}
	response := query(client, url)

	var tickerResponse TickerInfoResponse
	pTicker := &tickerResponse.Tickers

	if err := unmarshal(response, pTicker); err != nil {
		panic(err)
	}

	// print result
	for ticker, v := range tickerResponse.Tickers {
		fmt.Printf("%s High [%f]\t Avg [%f]\t Low[%f]\t Volume[%f]\t Current Volume[%f]\n",
			Bold(ticker), v.High, Green(v.Avg), v.Low, v.Vol, v.VolCur)
	}

}
