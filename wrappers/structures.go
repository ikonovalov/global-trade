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
	"fmt"
	"os"
)

type (
	CryptCurrencyExchange interface {
		GetTickers([]string, chan<- map[string]Ticker)
		GetBalances(ch chan<- Balance)
		Release()
	}

	Balance struct {
		Exchange       Exchange
		Funds          map[string]float64
		AvailableFunds map[string]float64
	}

	Balances []Balance

	Exchange struct {
		CryptCurrencyExchange
		Name  string
		SName string
		Link  string
		Cold  bool
	}

	ByExchangeName struct{ Balances }

	Ticker struct {
		High    float64
		Low     float64
		Avg     float64
		Vol     float64
		VolCur  float64
		Buy     float64
		Sell    float64
		Last    float64
		Updated int64
	}
)

func (s Balances) Len() int      { return len(s) }
func (s Balances) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByExchangeName) Less(i, j int) bool { return s.Balances[i].Exchange.Name < s.Balances[j].Exchange.Name }

func fatal(v ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprint(v))
	os.Exit(1)
}