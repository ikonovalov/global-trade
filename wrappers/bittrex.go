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

package wrappers

import (
	"github.com/toorop/go-bittrex"
	"time"
	"log"
	"fmt"
	"os"
	"github.com/ikonovalov/go-cloudflare-scraper"
	"net/http"
)

type BittrexWrapper struct {
	bittrex *bittrex.Bittrex
	availableMarkets map[string]bittrex.Market
}

type BittrexApiCredential struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func NewBittrex(credential BittrexApiCredential) *BittrexWrapper {
	cloudflare, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		fatal(err)
	}
	httpClient := &http.Client{Transport: cloudflare, Jar: cloudflare.Cookies, Timeout: time.Second * 60}
	bittrexClient := bittrex.NewWithCustomHttpClient(credential.Key, credential.Secret, httpClient)

	ba := BittrexWrapper{
		bittrex: bittrexClient,
		availableMarkets: make(map[string]bittrex.Market),
	}

	// upload markets
	start := time.Now()
	markets, err := bittrexClient.GetMarkets()
	elapsed := time.Since(start)
	log.Printf("Bittrex.GetMarkets took %s", elapsed)
	for _, m := range markets {
		ba.availableMarkets[m.MarketName] = m
	}

	return &ba
}

func (bw *BittrexWrapper) GetBalances( ch chan<- Balance) {
	start := time.Now()
	balances, err := bw.bittrex.GetBalances()
	elapsed := time.Since(start)
	log.Printf("Bittrex.GetBalances took %s", elapsed)
	if err != nil {
		fatal(err)
	}
	canonicalBalances := Balance{
		Exchange:       Exchange{Name: "Bittrex", Link: "https://bittrex.com"},
		Funds:          make(map[string]float64),
		AvailableFunds: make(map[string]float64),
	}
	for _, bb := range balances {
		balF64, _ := bb.Balance.Float64()
		avaF64, _ := bb.Available.Float64()
		canonicalBalances.Funds[bb.Currency] = balF64
		canonicalBalances.AvailableFunds[bb.Currency] = avaF64
	}
	ch <- canonicalBalances
}

func (bw *BittrexWrapper) GetTickers(paris []string, ch chan <- map[string]Ticker) {

}

func (bw *BittrexWrapper) Release()  {
	// nothing to do now
}

func fatal(v ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprint(v))
	os.Exit(1)
}
