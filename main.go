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
	"fmt"
	. "github.com/logrusorgru/aurora"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.")

	cmdInfo = app.Command("info", "Show all listed tickers on the Yobit").Default()

	cmdTicker     = app.Command("ticker", "Command provides statistic data for the last 24 hours.")
	cmdTickerPair = cmdTicker.Arg("pairs", "Listing ticker name. eth_btc, xem_usd, and so on.").Default("btc_usd").String()

	cmdDepth      = app.Command("depth", "Command returns information about lists of active orders for selected pairs.")
	cmdDepthPair  = cmdDepth.Arg("pairs", "eth_btc, xem_usd and so on.").Default("btc_usd").String()
	cmdDepthLimit = cmdDepth.Arg("limit", "Depth output limit").Default("20").Int()

	cmdTrades      = app.Command("trades", "Command returns information about the last transactions of selected pairs.")
	cmdTradesPair  = cmdTrades.Arg("pairs", "waves_btc, dash_usd and so on.").Default("btc_usd").String()
	cmdTradesLimit = cmdTrades.Arg("limit", "Trades output limit.").Default("100").Int()

	cmdBalances = app.Command("balances", "Command returns information about user's balances and priviledges of API-key as well as server time.")
)

func main() {

	kingpin.Version("0.1.0")

	yobit := NewYobit()

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
	case "balances":
		{
			channel := make(chan GetInfoResponse)
			go yobit.GetInfo(channel)
			getInfoRes := <-channel
			data := getInfoRes.Data
			printFunds("Balances (include orders)", data.FundsIncludeOrders, data.ServerTime)
		}
	default:
		panic("Unknown command " + command)
	}

}
func printFunds(caption string, funds map[string]float64, updated int64) {
	fmt.Printf("%s [%s]\n", Bold(caption), time.Unix(updated, 0).Format(time.Stamp))
	for k, v := range funds {
		if v == 0 {
			continue
		}
		fmt.Printf("%-5s %.8f\n", Bold(strings.ToUpper(k)), v)
	}
}

func printInfoRecord(desc PairInfo, name string) {
	Colored := Bold
	if desc.Hidden != 0 {
		Colored = Red
	}
	// TODO Try https://github.com/olekukonko/tablewriter
	fmt.Printf("%-12sfee %1.3f%% hidden %d Min amount %8.8f Min/Max price [%8.8f/%8.8f]\n",
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
