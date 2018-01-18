package main

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.")

	cmdInfo = app.Command("info", "Show all listed tickers on the Yobit")

	cmdTicker     = app.Command("ticker", "Client command: depth | ticker")
	cmdTickerPair = cmdTicker.Arg("pair", "Listing ticker name. eth_btc, xem_usd, and so on.").Default("btc_usd").String()

	cmdDepth      = app.Command("depth", "ASK/BID depth")
	cmdDepthPair  = cmdDepth.Arg("pair", "eth_btc, xem_usd and so on.").Default("btc_usd").String()
	cmdDepthLimit = cmdDepth.Arg("limit", "Depth limit").Default("20").Int()
)

func main() {

	kingpin.Version("0.1.0")

	yobit := Yobit{&http.Client{}}

	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch command {
	case "info":
		{
			channel := make(chan InfoResponse)
			go yobit.Info(channel)
			infoResponse := <-channel

			for name, desc := range infoResponse.Pairs {
				Colored := Bold
				if desc.Hidden != 0 {
					Colored = Red
				}
				fmt.Printf("%s\tfee %f hidden %d min_amount %f min_price %f max_price %f \n",
					Colored(Bold(strings.ToUpper(name))), desc.Fee, desc.Hidden, desc.MinAmount, desc.MinPrice, desc.MaxPrice)
			}
		}
	case "ticker":
		{
			channel := make(chan TickerInfoResponse)
			go yobit.Tickers24(*cmdTickerPair, channel)
			tickerResponse := <-channel

			for ticker, v := range tickerResponse.Tickers {
				printTicker(v, ticker)
			}
		}
	case "depth":
		{
			channel := make(chan DepthResponse)
			go yobit.DepthLimited(*cmdDepthPair, *cmdDepthLimit, channel)
			depthResponse := <-channel
			orders := depthResponse.Orders[*cmdDepthPair]

			for idx, ask := range orders.Asks {
				printDepth(idx, ask)
			}
		}
	default:
		panic("Unknown command " + command)
	}

}

func printDepth(idx int, ask Order) (int, error) {
	return fmt.Printf("#%d %.8f <- %.8f\n", idx+1, ask.Price, ask.Quantity)
}

func printTicker(v Ticker, tickerName string) {
	spread := v.Sell - v.Buy
	fmt.Printf(
		"%-9s H/A/L [%.8f/%.8f/%.8f] Buy[%.8f] Sell[%.8f] Last[%.8f] Spread[%.8f] Volume[%.8f] Current Volume[%.8f] \n",
		Bold(strings.ToUpper(tickerName)), v.High, Green(v.Avg), v.Low, v.Buy, v.Sell, Green(v.Last), BgRed(spread), v.Vol, v.VolCur)
}
