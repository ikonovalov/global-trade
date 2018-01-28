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
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
	"strconv"
)

var (
	bold []int = tablewriter.Colors{tablewriter.Bold}
	norm []int = tablewriter.Colors{0}
)

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

func printWallets(coinFilter string, fundsAndTickers struct {
	funds   map[string]float64
	tickers map[string]Ticker
}, updated int64) {
	// setup table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Coin", "Hold", "USD RATE (24AVG)", "USD CAP (24AVG)", "USD CAP (LAST)", "DIFF (ABS)", "DIFF(%)"})
	table.SetHeaderColor(bold, bold, bold, bold, bold, bold, bold)
	table.SetColumnColor(bold, norm, norm, norm, norm, norm, norm)

	// determinate price multiplication indicator
	basePriceFunc := func(ticker Ticker) float64 { return ticker.Avg }
	actualPriceFunc := func(ticker Ticker) float64 { return ticker.Last }

	var (
		baseUsdTotal   float64
		actualUsdTotal float64
		diffUsdTotal   float64
	)

	for coin, volume := range fundsAndTickers.funds {
		if volume == 0 || (coinFilter != "all" && coin != coinFilter) {
			continue
		}
		tickerName := fmt.Sprintf("%s_usd", coin)

		basePrice := basePriceFunc(fundsAndTickers.tickers[tickerName])
		baseUsdCoinPrice := volume * basePrice
		baseUsdTotal += baseUsdCoinPrice

		actualPrice := actualPriceFunc(fundsAndTickers.tickers[tickerName])
		actualUsdCoinPrice := volume * actualPrice
		actualUsdTotal += actualUsdCoinPrice

		diffUsdCoinPriceAbs := actualUsdCoinPrice - baseUsdCoinPrice
		diffUsdTotal += diffUsdCoinPriceAbs

		var diffUsdCoinPricePercent float64
		if diffUsdCoinPriceAbs != 0 {
			diffUsdCoinPricePercent = diffUsdCoinPriceAbs / baseUsdCoinPrice * float64(100)
		}

		usdCoinColor := Green
		if basePrice > actualPrice {
			usdCoinColor = Red
		} else if basePrice == actualPrice {
			usdCoinColor = Gray
		}

		table.Append([]string{
			strings.ToUpper(coin),
			fmt.Sprintf("%.8f", volume),
			fmt.Sprintf("%.8f", basePrice),
			fmt.Sprintf("%.8f", baseUsdCoinPrice),
			fmt.Sprintf("%.8f", usdCoinColor(actualUsdCoinPrice)),
			fmt.Sprintf("%+8.8f", usdCoinColor(diffUsdCoinPriceAbs)),
			fmt.Sprintf("%+3.2f", usdCoinColor(diffUsdCoinPricePercent)),
		})
	}
	table.SetFooter([]string{
		"",
		time.Unix(updated, 0).Format(time.Stamp),
		"Total cap",
		fmt.Sprintf("%8.2f", baseUsdTotal),
		fmt.Sprintf("%8.2f", actualUsdTotal),
		fmt.Sprintf("%8.2f", diffUsdTotal),
		"",
	})
	table.Render()
}

func printDepth(idx int, ask Offer) {
	fmt.Printf("#%d %.8f <- %.8f\n", idx+1, ask.Price, ask.Quantity)
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

		fmt.Printf("%s %s Price[%.8f] Amount[%.8f] \u21D0 %d\n", tm, Bold(Colored(tradeDirection)), trade.Price, trade.Amount, trade.Tid)
	}
}

func printActiveOrders(activeOrders ActiveOrdersResponse) {
	for ordId, ord := range activeOrders.Orders {
		created, _ := strconv.ParseInt(ord.Created, 10, 64)
		fmt.Printf("%s ID[%s] %s amount: %.8f rate: %.8f\n",
			time.Unix(created, 0).Format(time.Stamp), ordId, strings.ToUpper(ord.Type), ord.Amount, ord.Rate)
	}
}
