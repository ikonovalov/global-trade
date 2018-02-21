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
	"net/http"
	"time"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/shopspring/decimal"
	"log"
)

var (
	EtherScan = Exchange{Name: "Ethereum", Link: "etherscan.io", Cold: true}
	client    = http.Client{Timeout: time.Second * 10}
)

type (
	EtherScanAccountBalancesResponse struct {
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Result  []EthereumBalance `json:"result"`
	}

	EthereumBalance struct {
		Account string          `json:"account"`
		Balance decimal.Decimal `json:"balance"`
	}

	EthereumBalances []EthereumBalance
)

func (e *EthereumBalance) ToBalance() Balance {
	balanceF64wei, exec := e.Balance.Float64()
	if !exec {
		fatal("Decimal balance overflow")
	}
	funds := map[string]float64{"ETH": balanceF64wei / 1000000000000000000.0}
	return Balance{
		Exchange:       EtherScan,
		Funds:          funds,
		AvailableFunds: funds,
	}
}

func (e EthereumBalances) SummaryBalance() Balance {
	totalBalance := 0.0
	for _, b := range e {
		balanceF64, exec := b.Balance.Float64()
		if !exec {
			fatal("Decimal overflow for ", b.Account)
		}
		totalBalance += balanceF64 / 1000000000000000000.0
	}
	funds := map[string]float64{"ETH": totalBalance}

	return Balance{
		Exchange:       EtherScan,
		Funds:          funds,
		AvailableFunds: funds,
	}
}

type EtherScanCredential struct {
	Accounts []string `json:"accounts"`
}

func GetEthereumBalances(addresses []string, ch chan<- EthereumBalances) {
	addressesLine := strings.Join(addresses, ",")
	queryString := fmt.Sprintf(
		"https://api.etherscan.io/api?module=%s&action=%s&address=%s&tag=latest",
		"account",
		"balancemulti",
		addressesLine,
	)
	start := time.Now()
	resp, err := client.Get(queryString)
	elapsed := time.Since(start)
	log.Printf("EtherScan.Account.BalanceMulti took %s", elapsed)
	if err != nil {
		fatal(err)
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	var responseStructure EtherScanAccountBalancesResponse
	err = json.Unmarshal(responseBytes, &responseStructure)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	ch <- responseStructure.Result
}
