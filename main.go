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
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

var (
	app = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.")

	cmdMarkets      = app.Command("markets", "Show all listed tickers on the Yobit").Default()
	cmdInfoCurrency = cmdMarkets.Arg("cryptocurrency", "Show markets only for specified currency: btc, eth, usd and so on.").Default("").String()

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
	case "markets":
		{
			channel := make(chan InfoResponse)
			go yobit.Info(channel)
			infoResponse := <-channel
			printInfoRecords(infoResponse, *cmdInfoCurrency)
			fmt.Printf("\nTotal markets %d\n", len(infoResponse.Pairs))
		}
	case "ticker":
		{
			channel := make(chan TickerInfoResponse)
			go yobit.Tickers24(strings.ToLower(*cmdTickerPair), channel)
			tickerResponse := <-channel

			for ticker, v := range tickerResponse.Tickers {
				printTicker(v, ticker)
			}
		}
	case "depth":
		{
			channel := make(chan DepthResponse)
			go yobit.DepthLimited(strings.ToLower(*cmdDepthPair), *cmdDepthLimit, channel)
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
			go yobit.TradesLimited(strings.ToLower(*cmdTradesPair), *cmdTradesLimit, channel)
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
func printInfoRecords(infoResponse InfoResponse, currencyFilter string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Market", "Hidden", "Min amount", "Min price", "Max price"})
	bold := tablewriter.Colors{tablewriter.Bold}
	norm := tablewriter.Colors{0}
	table.SetHeaderColor(bold, bold, bold, bold, bold)
	table.SetColumnColor(bold, norm, norm, norm, norm)

	currencyFilter = strings.ToUpper(currencyFilter)
	for name, desc := range infoResponse.Pairs {
		hidden := "NO"
		if desc.Hidden == 1 {
			hidden = "YES"
		}
		if currencyName := strings.ToUpper(name); currencyFilter == "" || strings.Contains(currencyName, currencyFilter) {
			table.Append([]string{
				currencyName,
				fmt.Sprintf("%s", hidden),
				fmt.Sprintf("%f", desc.MinAmount),
				fmt.Sprintf("%f", desc.MinPrice),
				fmt.Sprintf("%f", desc.MaxPrice),
			})
		}
	}
	table.Render()
}

func printFunds(caption string, funds map[string]float64, updated int64) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Coin", "Hold"})
	bold := tablewriter.Colors{tablewriter.Bold}
	norm := tablewriter.Colors{0}
	table.SetHeaderColor(bold, bold)
	table.SetColumnColor(bold, norm)
	fmt.Printf("%s [%s]\n", Bold(caption), time.Unix(updated, 0).Format(time.Stamp))
	for k, v := range funds {
		if v == 0 {
			continue
		}
		table.Append([]string{strings.ToUpper(k), fmt.Sprintf("%.8f", v)})
	}
	table.Render()
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
