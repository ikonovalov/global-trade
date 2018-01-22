package main

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
	"github.com/ikonovalov/go-cloudflare-scraper"
)

var (
	app = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.")

	cmdInfo = app.Command("info", "Show all listed tickers on the Yobit")

	cmdTicker     = app.Command("ticker", "Client command: depth | ticker")
	cmdTickerPair = cmdTicker.Arg("pairs", "Listing ticker name. eth_btc, xem_usd, and so on.").Default("btc_usd").String()

	cmdDepth      = app.Command("depth", "ASK/BID depth")
	cmdDepthPair  = cmdDepth.Arg("pairs", "eth_btc, xem_usd and so on.").Default("btc_usd").String()
	cmdDepthLimit = cmdDepth.Arg("limit", "Depth output limit").Default("20").Int()

	cmdTrades      = app.Command("trades", "Last trades information for pairs")
	cmdTradesPair  = cmdTrades.Arg("pairs", "waves_btc, dash_usd and so on.").Default("btc_usd").String()
	cmdTradesLimit = cmdTrades.Arg("limit", "Trades output limit.").Default("100").Int()
)

func main() {

	kingpin.Version("0.1.0")

	scraper, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		panic(err)
	}

	keys, err := loadApiKeys()
	if err != nil {
		panic(err)
	}

	yobit := Yobit{
		client: &http.Client{
			Transport: scraper,
			Jar:       scraper.Cookies,
		},
		apiKeys: &keys,
	}

	yobit.GetInfo()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch command {
	case "info":
		{
			channel := make(chan InfoResponse)
			go yobit.Info(channel)
			infoResponse := <-channel

			for name, desc := range infoResponse.Pairs {
				printInfoRecord(desc, name)
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
			fmt.Println(Bold(strings.ToUpper(*cmdDepthPair)))
			for idx, ask := range orders.Asks {
				printDepth(idx, ask)
			}
		}
	case "trades":
		{
			channel := make(chan TradesResponse)
			go yobit.TradesLimited(*cmdTradesPair, *cmdTradesLimit, channel)
			tradesResponse := <-channel
			for ticker, trades := range tradesResponse.Trades {
				fmt.Println(Bold(strings.ToUpper(ticker)))
				printTrades(trades)

			}
		}
	default:
		panic("Unknown command " + command)
	}

}
func printInfoRecord(desc PairInfo, name string) {
	Colored := Bold
	if desc.Hidden != 0 {
		Colored = Red
	}
	// TODO Try https://github.com/olekukonko/tablewriter
	fmt.Printf("%s\tfee %.3f%% hidden %d min_amount %f min_price %f max_price %f \n",
		Colored(Bold(strings.ToUpper(name))), desc.Fee, desc.Hidden, desc.MinAmount, desc.MinPrice, desc.MaxPrice)
}

func printDepth(idx int, ask Order) (int, error) {
	return fmt.Printf("#%d %.8f <- %.8f\n", idx+1, ask.Price, ask.Quantity)
}

func printTicker(v Ticker, tickerName string) {
	spread := v.Sell - v.Buy
	updated := time.Unix(v.Updated, 0).Format(time.Stamp)
	fmt.Printf(
		"%s %-9s H/A/L [%.8f/%.8f/%.8f] Buy[%.8f] Sell[%.8f] Last[%.8f] Spread[%.8f] Volume[%.8f] Current Volume[%.8f] \n",
		updated, Bold(strings.ToUpper(tickerName)), v.High, Green(v.Avg), v.Low, v.Buy, v.Sell, Green(v.Last), BgRed(spread), v.Vol, v.VolCur)
}

func printTrades(trades []Trade) {
	for _, trade := range trades {
		tm := time.Unix(trade.Timestamp, 0).Format(time.Stamp)
		Colored := BgGreen
		tradeDirection := "Buy "
		if trade.Type == "ask" {
			Colored = BgRed
			tradeDirection = "Sell"
		}

		fmt.Printf("%s %s Price[%.8f] Amount[%.8f] \u21D0 %d\n", tm, Colored(tradeDirection), trade.Price, trade.Amount, trade.Tid)
	}
}
