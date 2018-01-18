package main

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	kingpin.Version("0.0.1")
	kingpin.Parse()

	client := &http.Client{}
	ch := make(chan TickerInfoResponse)
	go tickers24(client, *ticker, ch)
	tickerResponse := <-ch

	// print result
	for ticker, v := range tickerResponse.Tickers {
		spread := v.Sell - v.Buy
		fmt.Printf(
			"%-9s H/A/L [%.8f/%.8f/%.8f] Buy[%.8f] Sell[%.8f] Last[%.8f] Spread[%.8f] Volume[%.8f] Current Volume[%.8f] \n",
			Bold(strings.ToUpper(ticker)), v.High, Green(v.Avg), v.Low, v.Buy, v.Sell, Green(v.Last), BgRed(spread), v.Vol, v.VolCur)
	}

	ch2 := make(chan DepthResponse)
	go depth(client, *ticker, ch2)
	depthResponse := <- ch2
	orders := depthResponse.Orders[*ticker]

	for idx, ask := range orders.Asks {
		fmt.Printf("#%d %.8f <- %.8f\n" , idx + 1, ask.Price, ask.Quantity)
	}

}


