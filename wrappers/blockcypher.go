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
	"github.com/blockcypher/gobcy"
	"time"
	"math"
)

var(
	BlockCypher = Exchange{Name: "LiteCoin", Link: "blockcypher.com", Cold: true}
	litecoinDecimalPoint = math.Pow(10.0, 8.0)
)

type (
	BlockCypherCredential struct {
		LTC []string `json:"ltc,omitempty"`
	}
	BlochCypherBalances []gobcy.Addr
)

func (bcb BlochCypherBalances) SummaryBalance() (Balance) {
	totalBalance := 0.0
	for _, addr := range bcb {
		totalBalance += float64(addr.FinalBalance)
	}

	funds := map[string]float64{ "LTC": totalBalance / litecoinDecimalPoint}

	return Balance{
		Exchange:       BlockCypher,
		Funds: funds,
		AvailableFunds: funds,
	}
}

func GetLiteCoinBalances(accounts []string, ch chan <- BlochCypherBalances)  {
	btc := gobcy.API{ Coin:"ltc", Chain: "main"}
	accLen := len(accounts)
	rs := make([]gobcy.Addr, 0, accLen)
	for _, acc := range accounts {
		addr, err := btc.GetAddrBal(acc, nil)
		if err != nil {
			fatal(err)
		}
		rs = append(rs, addr)
		if accLen > 1 {
			time.Sleep(time.Millisecond * 200)
		}
	}
	ch <- rs
}
