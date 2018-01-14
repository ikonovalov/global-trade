package main

import (
	"encoding/json"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"net/http"
)

// this is a comment

const (
	API_BASE    = "https://yobit.net/api/"
	API_VERSION = "3"
)

func main() {

	url := "https://yobit.net/api/3/ticker/xem_usd-eth_usd-xem_eth"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Errorf("error. HTTP %d", resp.StatusCode)
	} else {
		var tickerResponse TickerInfoResponse
		if err := json.Unmarshal(response, &tickerResponse.Tickers); err != nil {
			panic(err)
		}
		for ticker, v := range tickerResponse.Tickers {
			fmt.Printf("%s High [%f]\t Low[%f]\t Volume[%f]\t Current Volume[%f]\n",
				Bold(ticker), v.High, v.Low, v.Vol, v.VolCur)
		}

	}
}
