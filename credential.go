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
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/ikonovalov/go-yobit"
	"github.com/ikonovalov/global-trade/wrappers"
)

type GlobalCredentials struct {
	Version    uint16                        `json:"version"`
	Encryption bool                          `json:"encryption"`
	Yobit      wrappers.YobitApiCredential   `json:"yobit,omitempty"`
	Bittrex    wrappers.BittrexApiCredential `json:"bittrex,omitempty"`
}

func loadApiCredential() (GlobalCredentials, error) {
	file, e := ioutil.ReadFile(credentialFile)
	if e != nil {
		return GlobalCredentials{}, e
	}
	var keys GlobalCredentials
	unmarshalError := json.Unmarshal(file, &keys)

	return keys, unmarshalError
}

func createCredentialFile(adiCredential yobit.ApiCredential) {
	if _, err := os.Stat(credentialFile); os.IsNotExist(err) {
		if _, err = os.Create(credentialFile); err != nil {
			panic(err)
		}
	}
	data, _ := json.Marshal(adiCredential)
	if err := ioutil.WriteFile(credentialFile, data, 0644); err != nil {
		panic(err)
	}

}
