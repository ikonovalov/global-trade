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

package bittrex_async

import (
	"github.com/toorop/go-bittrex"
	"fmt"
	"os"
	"time"
	"net/http"
	"github.com/ikonovalov/go-cloudflare-scraper"
	"log"
)

type BittrexAsync struct {
	*bittrex.Bittrex
}

func New(credential ApiCredential) *BittrexAsync {
	cloudflare, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		fatal(err)
	}
	httpClient := &http.Client{Transport: cloudflare, Jar: cloudflare.Cookies, Timeout: time.Second * 60}
	ba := BittrexAsync{
		Bittrex: bittrex.NewWithCustomHttpClient(
			credential.Key, credential.Secret,
			httpClient,
		),
	}
	return &ba
}

func (ba *BittrexAsync) MarketsAsync(ch chan<- []bittrex.Market) {
	markets, err := ba.GetMarkets()
	if err != nil {
		fatal(err)
	}
	ch <- markets
}

func (ba *BittrexAsync) GetBalancesAsync(ch chan []bittrex.Balance) {
	start := time.Now()
	balances, err := ba.GetBalances()
	elapsed := time.Since(start)
	log.Printf("Bittrex.GetBalances took %s", elapsed)
	if err != nil {
		fatal(err)
	}
	ch <- balances
}

func (ba *BittrexAsync) GetTickers(ch chan<- []bittrex.Ticker, markets ...string) {
	if len(markets) == 0 {
		ch <- []bittrex.Ticker{}
		return
	}
	tickerRs := make([]bittrex.Ticker, 0, len(markets))
	for _, m := range markets {
		ticker, err := ba.GetTicker(m)
		if err != nil {
			fatal(err)
		}
		tickerRs = append(tickerRs, ticker)
	}
	ch <- tickerRs
	return
}

func fatal(v ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprint(v))
	os.Exit(1)
}
