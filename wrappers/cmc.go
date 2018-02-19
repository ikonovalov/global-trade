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
	coinApi "github.com/miguelmota/go-coinmarketcap"
	"time"
	"log"
)

type CoinMarketCap struct {
}

func (mc *CoinMarketCap) GetMarketData(ch chan<- map[string]coinApi.Coin) {
	start := time.Now()
	top, err := coinApi.GetAllCoinData(1000)
	if err != nil {
		fatal(err)
	}
	elapsed := time.Since(start)
	log.Printf("CMC.GetMarketData (TOP100) took %s", elapsed)

	rs := make(map[string]coinApi.Coin)
	for _, coin := range top {
		rs[coin.Symbol] = coin
	}
	ch <- rs
}
