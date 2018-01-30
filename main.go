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
	"log"
	"io/ioutil"
)

var (
	app            = kingpin.New("yobit", "Yobit cryptocurrency exchange crafted client.").Version("0.2.0")
	appVerboseFlag = app.Flag("verbose", "Print additional information").Bool()

	cmdInit = app.Command("init", "Initialize nonce and keys container")

	cmdMarkets      = app.Command("markets", "(m) Show all listed tickers on the Yobit").Alias("m")
	cmdInfoCurrency = cmdMarkets.Arg("cryptocurrency", "Show markets only for specified currency: btc, eth, usd and so on.").Default("").String()

	cmdTicker     = app.Command("ticker", "(tc) Command provides statistic data for the last 24 hours.").Alias("tc")
	cmdTickerPair = cmdTicker.Arg("pairs", "Listing ticker name. eth_btc, xem_usd, and so on.").Default("btc_usd").String()

	cmdDepth         = app.Command("depth", "(d) Command returns information about lists of active orders for selected pairs.").Alias("d")
	cmdDepthAskOrBid = cmdDepth.Arg("direction", "ask or bid").Default("ask").String()
	cmdDepthPair     = cmdDepth.Arg("pairs", "eth_btc, xem_usd and so on.").Default("btc_usd").String()
	cmdDepthLimit    = cmdDepth.Arg("limit", "Depth output limit").Default("20").Int()

	cmdTrades      = app.Command("trades", "(tr) Command returns information about the last transactions of selected pairs.").Alias("tr")
	cmdTradesPair  = cmdTrades.Arg("pairs", "waves_btc, dash_usd and so on.").Default("btc_usd").String()
	cmdTradesLimit = cmdTrades.Arg("limit", "Trades output limit.").Default("100").Int()

	cmdWallets             = app.Command("wallets", "(w) Command returns information about user's balances and privileges of API-key as well as server time.").Alias("w")
	cmdWalletsBaseCurrency = cmdWallets.Arg("base", "Base recalculated currency").Default("usd").String()

	cmdActiveOrders    = app.Command("active-orders", "(ao) Show active orders").Alias("ao")
	cmdActiveOrderPair = cmdActiveOrders.Arg("pair", "doge_usd...").Required().String()

	cmdTradeHistory     = app.Command("trade-history", "(th) Trade history").Alias("th")
	cmdTradeHistoryPair = cmdTradeHistory.Arg("pair", "doge_usd...").Required().String()

	cmdTrade       = app.Command("trade", "(t) Creating new orders for stock exchange trading").Alias("t")
	cmdTradePair   = cmdTrade.Arg("pair", "Pair").Required().String()
	cmdTradeType   = cmdTrade.Arg("type", "Transaction type: sell or buy").Required().String()
	cmdTradeRate   = cmdTrade.Arg("rate", "Exchange rate for buying or selling").Required().Float64()
	cmdTradeAmount = cmdTrade.Arg("amount", "Exchange rate for buying or selling").Required().Float64()

	cmdBuy       = app.Command("buy", "(b) Buy on stock exchange").Alias("b")
	cmdBuyPair   = cmdBuy.Arg("pair", "Pair").Required().String()
	cmdBuyRate   = cmdBuy.Arg("rate", "Exchange rate for buying or selling").Required().Float64()
	cmdBuyAmount = cmdBuy.Arg("amount", "Exchange rate for buying or selling").Required().Float64()

	cmdSell       = app.Command("sell", "(s) Sell on stock exchange").Alias("s")
	cmdSellPair   = cmdSell.Arg("pair", "Pair").Required().String()
	cmdSellRate   = cmdSell.Arg("rate", "Exchange rate for buying or selling").Required().Float64()
	cmdSellAmount = cmdSell.Arg("amount", "Exchange rate for buying or selling").Required().Float64()

	cmdCancelOrder        = app.Command("cancel", "(c) Cancels the chosen order").Alias("c")
	cmdCancelOrderOrderId = cmdCancelOrder.Arg("order_id", "Order ID").Required().String()
)

func main() {

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// setup logging
	if !*appVerboseFlag {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	yobit := NewYobit()

	switch command {
	case "init":
		{
			createKeysFile()
		}
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
			go yobit.Tickers24([]string{strings.ToLower(*cmdTickerPair)}, channel)
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
			fmt.Printf("%s %s\n", Bold(strings.ToUpper(*cmdDepthPair)).String(), Bold(strings.ToUpper(*cmdDepthAskOrBid)).String())
			offers := &orders.Asks
			if strings.ToLower(*cmdDepthAskOrBid) == "bid" {
				offers = &orders.Bids
			}
			printDepth(offers)

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
			funds := data.FundsIncludeOrders
			usdPairs := make([]string, 0, len(funds))
			for coin, volume := range funds {
				pair := fmt.Sprintf("%s_%s", coin, *cmdWalletsBaseCurrency)
				if volume > 0 && yobit.isMarketExists(pair) {
					usdPairs = append(usdPairs, pair)
				}
			}
			tickersChan := make(chan TickerInfoResponse)
			go yobit.Tickers24(usdPairs, tickersChan)
			tickerRs := <-tickersChan
			fundsAndTickers := struct {
				funds   map[string]float64
				tickers map[string]Ticker
			}{data.FundsIncludeOrders, tickerRs.Tickers}
			printWallets(*cmdWalletsBaseCurrency, fundsAndTickers, data.ServerTime)
		}
	case "active-orders":
		{
			channel := make(chan ActiveOrdersResponse)
			go yobit.ActiveOrders(*cmdActiveOrderPair, channel)
			activeOrders := <-channel
			printActiveOrders(activeOrders)

		}
	case "trade-history":
		{
			channel := make(chan TradeHistoryResponse)
			go yobit.TradeHistory(*cmdTradeHistoryPair, channel)
			history := <-channel
			fmt.Println(history)
		}
	case "trade":
		{
			channel := make(chan TradeResponse)
			go yobit.Trade(*cmdTradePair, *cmdTradeType, *cmdTradeRate, *cmdTradeAmount, channel)
			trade := <-channel
			fmt.Printf("Order %d created\n", trade.Result.OrderId)
		}
	case "buy":
		{
			channel := make(chan TradeResponse)
			go yobit.Trade(*cmdBuyPair, "buy", *cmdBuyRate, *cmdBuyAmount, channel)
			trade := <-channel
			fmt.Printf("Order %d created\n", trade.Result.OrderId)
		}
	case "sell":
		{
			channel := make(chan TradeResponse)
			go yobit.Trade(*cmdSellPair, "sell", *cmdSellRate, *cmdSellAmount, channel)
			trade := <-channel
			fmt.Printf("Order %d created\n", trade.Result.OrderId)
		}
	case "cancel":
		{
			channel := make(chan CancelOrderResponse)
			go yobit.CancelOrder(*cmdCancelOrderOrderId, channel)
			cancelResult := <-channel
			fmt.Printf("Order %d candeled\n", cancelResult.Result.OrderId)
		}
	default:
		panic("Unknown command " + command)
	}

}
