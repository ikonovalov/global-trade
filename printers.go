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
