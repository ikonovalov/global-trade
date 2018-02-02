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
	"math"
)

var (
	bold []int = tablewriter.Colors{tablewriter.Bold}
	norm []int = tablewriter.Colors{0}
)

func sprintf64(v float64) string {
	return fmt.Sprintf("%8.8f", v)
}

func fatal(v ...interface{}) {
	fmt.Println(Red(Bold(fmt.Sprint(v))).String())
	os.Exit(1)
}

func printInfoRecords(infoResponse InfoResponse, currencyFilter string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Market", "Hidden", "Fee", "Min amount", "Min price", "Max price"})
	bold := tablewriter.Colors{tablewriter.Bold}
	norm := tablewriter.Colors{0}
	table.SetHeaderColor(bold, bold, bold, bold, bold, bold)
	table.SetColumnColor(bold, norm, norm, norm, norm, norm)

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
				fmt.Sprintf("%2.2f%%", desc.Fee),
				fmt.Sprintf("%8.8f", desc.MinAmount),
				fmt.Sprintf("%8.8f", desc.MinPrice),
				fmt.Sprintf("%8.8f", desc.MaxPrice),
			})
		}
	}
	table.Render()
}

func printWallets(baseCurrency string, fundsAndTickers struct {
	funds   map[string]float64
	tickers map[string]Ticker
}, updated int64) {
	// setup table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Coin",
		"Hold",
		fmt.Sprintf("%s RATE (24AVG)", baseCurrency),
		fmt.Sprintf("%s CAP (24AVG)", baseCurrency),
		fmt.Sprintf("%s CAP (LAST)", baseCurrency),
		"DIFF (ABS)",
		"DIFF (%)",
	})
	table.SetHeaderColor(bold, bold, bold, bold, bold, bold, bold)
	table.SetColumnColor(bold, norm, norm, norm, norm, norm, norm)

	// determinate price multiplication indicator
	basePriceFunc := func(tickerName string) float64 {
		if tickerName == fmt.Sprintf("%s_%[1]s", baseCurrency) {
			return float64(1)
		} else {
			return fundsAndTickers.tickers[tickerName].Avg
		}
	}
	actualPriceFunc := func(tickerName string) float64 {
		if tickerName == fmt.Sprintf("%s_%[1]s", baseCurrency) {
			return float64(1)
		} else {
			return fundsAndTickers.tickers[tickerName].Last
		}
	}

	var (
		baseUsdTotal   float64
		actualUsdTotal float64
		diffUsdTotal   float64
	)

	for coin, volume := range fundsAndTickers.funds {
		if volume == 0 {
			continue
		}
		tickerName := fmt.Sprintf("%s_%s", coin, baseCurrency)

		basePrice := basePriceFunc(tickerName)
		baseUsdCoinPrice := volume * basePrice
		baseUsdTotal += baseUsdCoinPrice

		actualPrice := actualPriceFunc(tickerName)
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
		fmt.Sprintf("%+8.2f", diffUsdTotal),
		"",
	})
	table.Render()
}

func printOffers(offers Offers) {
	var (
		asks    = offers.Asks
		bids    = offers.Bids
		asksLen = len(asks)
		bidsLen = len(bids)
		depth   = math.Max(float64(asksLen), float64(bidsLen))
	)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"#",
		"ask price",
		"ask quantity",
		"bid price",
		"bid quantity",
	})
	table.SetHeaderColor(bold, bold, bold, bold, bold)
	table.SetColumnColor(norm, norm, norm, norm, norm)

	appendOffer := func(row []string, offer Offer) []string {
		return append(row, sprintf64(offer.Price), sprintf64(offer.Quantity))
	}

	appendEmpty := func(row []string) []string {
		return append(row, "", "")
	}

	for i := 0; i < int(depth); i++ {
		row := append(make([]string, 0, 4), fmt.Sprintf("%d", i+1))

		if i < asksLen {
			ask := asks[i]
			row = appendOffer(row, ask)
		} else {
			row = appendEmpty(row)
		}
		if i < bidsLen {
			bid := bids[i]
			row = appendOffer(row, bid)
		} else {
			row = appendEmpty(row)
		}
		table.Append(row)
	}
	table.Render()
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

func printTradeResult(trade TradeResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"OrderId",
		"Received",
		"Remains",
	})
	table.SetHeaderColor(bold, bold, bold)
	table.SetColumnColor(bold, norm, norm)
	table.Append([]string{
		fmt.Sprintf("%d", trade.OrderId),
		fmt.Sprintf("%8.8f", trade.Received),
		fmt.Sprintf("%8.8f", trade.Remains),
	})
	table.Render()
}

func printActiveOrders(activeOrders ActiveOrdersResponse) {
	for ordId, ord := range activeOrders.Orders {
		created, _ := strconv.ParseInt(ord.Created, 10, 64)
		fmt.Printf("%s ID[%s] %s amount: %.8f rate: %.8f\n",
			time.Unix(created, 0).Format(time.Stamp), ordId, strings.ToUpper(ord.Type), ord.Amount, ord.Rate)
	}
}
