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
	"github.com/ikonovalov/go-yobit"
)

type YobitWrapper struct {
	yobit *yobit.Yobit
}

type YobitApiCredential struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func NewYobit(credential YobitApiCredential) *YobitWrapper {
	yobt := yobit.New(yobit.ApiCredential{
		Key: credential.Key,
		Secret: credential.Secret,
	})
	return &YobitWrapper{yobit:yobt}
}

func (yw *YobitWrapper) Release() {
	yw.yobit.Release()
}

func (yw *YobitWrapper) GetBalances(ch chan<- Balances) {
	channelYobit := make(chan yobit.GetInfoResponse)
	go yw.yobit.GetInfo(channelYobit)
	yobitGetInfoRes := <-channelYobit
	data := yobitGetInfoRes.Data

	yobitBalances := Balances{
		Exchange:       Exchange{Name: "Yobit", Link: yobit.Url},
		Funds:          data.FundsIncludeOrders,
		AvailableFunds: data.Funds,
		Tickers:        make(map[string]Ticker),
	}
	ch <- yobitBalances
}

func tickerMapFromYobit(ytm map[string]yobit.Ticker) map[string]Ticker {
	rs := make(map[string]Ticker)
	for k,yt := range ytm {
		rs[k] = Ticker {
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
	return rs
}
