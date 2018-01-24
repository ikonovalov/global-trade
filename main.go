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
)

var (
	app = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.").Version("0.1.1")

	cmdMarkets      = app.Command("markets", "Show all listed tickers on the Yobit").Alias("m")
	cmdInfoCurrency = cmdMarkets.Arg("cryptocurrency", "Show markets only for specified currency: btc, eth, usd and so on.").Default("").String()

	cmdTicker     = app.Command("ticker", "Command provides statistic data for the last 24 hours.").Alias("tc")
	cmdTickerPair = cmdTicker.Arg("pairs", "Listing ticker name. eth_btc, xem_usd, and so on.").Default("btc_usd").String()

	cmdDepth      = app.Command("depth", "Command returns information about lists of active orders for selected pairs.").Alias("dp")
	cmdDepthPair  = cmdDepth.Arg("pairs", "eth_btc, xem_usd and so on.").Default("btc_usd").String()
	cmdDepthLimit = cmdDepth.Arg("limit", "Depth output limit").Default("20").Int()

	cmdTrades      = app.Command("trades", "Command returns information about the last transactions of selected pairs.").Alias("tr")
	cmdTradesPair  = cmdTrades.Arg("pairs", "waves_btc, dash_usd and so on.").Default("btc_usd").String()
	cmdTradesLimit = cmdTrades.Arg("limit", "Trades output limit.").Default("100").Int()

	cmdWallets = app.Command("wallets", "Command returns information about user's balances and priviledges of API-key as well as server time.").Alias("bl")

	cmdActiveOrders = app.Command("active-orders", "Show active orders").Alias("ao")
	cmdActiveOrderPair = cmdActiveOrders.Arg("pair", "doge_usd...").Required().String()
)

func main() {

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
			orders := depthResponse.Offers[*cmdDepthPair]
			fmt.Println(Bold(strings.ToUpper(*cmdDepthPair)))
			fmt.Println(Bold("ASK"))
			for idx, ask := range orders.Asks {
				printDepth(idx, ask) // TODO Prints BIDS and add table output
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
	case "wallets":
		{
			channel := make(chan GetInfoResponse)
			go yobit.GetInfo(channel)
			getInfoRes := <-channel
			data := getInfoRes.Data
			printFunds("Balances (include orders)", data.FundsIncludeOrders, data.ServerTime)
		}
	case "active-orders":
		{
			channel := make(chan ActiveOrdersResponse)
			go yobit.ActiveOrders("xem_eth", channel)
			activeOrders := <-channel
			fmt.Println(activeOrders)
		}
	default:
		panic("Unknown command " + command)
	}

}
