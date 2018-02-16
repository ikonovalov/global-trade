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
	"github.com/ikonovalov/go-yobit"
	w "github.com/ikonovalov/global-trade/wrappers"
)

// Yobit converters ================================================

func tickerFromYobit(yt yobit.Ticker) w.Ticker {
	return w.Ticker{
		High: yt.High,
		Low: yt.Low,
		Avg: yt.Avg,
		Vol: yt.Vol,
		VolCur: yt.VolCur,
		Buy: yt.Buy,
		Sell: yt.Sell,
		Last: yt.Last,
		Updated: yt.Updated,
	}
}

func tickerMapFromYobit(ytm map[string]yobit.Ticker) map[string]w.Ticker {
	rs := make(map[string]w.Ticker)
	for k,v := range ytm {
		rs[k] = tickerFromYobit(v)
	}
	return rs
}

// Bittrex converters ===============================================