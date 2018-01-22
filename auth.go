package main

import (
	"os"
	"io/ioutil"
	"strconv"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
)

type ApiKeys struct {
	Key string `json:"key"`
	Secret string `json:"secret"`
}

func (y *Yobit) GetNonce() (nonce uint64) {
	createNonceFileIfNotExists()
	return readNonce()
}

func readNonce() uint64 {
	data, e := ioutil.ReadFile("nonce")
	if e != nil {
		panic("Nonce file read error. ")
	}
	nonce, conErr := strconv.ParseUint(string(data), 10, 64)
	if conErr != nil {
		panic(conErr)
	}
	return nonce
}

func writeNonce(data []byte) {
	if err := ioutil.WriteFile("nonce", data, 0644); err != nil {
		panic(err)
	}
}

func createNonceFileIfNotExists() {
	if _, err := os.Stat("nonce"); os.IsNotExist(err) {
		if _, err = os.Create("nonce"); err != nil {
			panic(err)
		}
		d1 := []byte("1")
		writeNonce(d1)
	}
}

func incrementNonce(nonceOld *uint64) {
	newNonce := *nonceOld + 1
	ns := strconv.FormatUint(newNonce, 10)
	writeNonce([]byte(ns))
}

func signHmacSha512(secret []byte, message []byte) (digest string) {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	digest = hex.EncodeToString(mac.Sum(nil))
	return
}

func loadApiKeys() (ApiKeys, error) {
	file, e := ioutil.ReadFile("keys")
	if e != nil {
		panic(e)
	}
	var keys ApiKeys
	unmarshalError := json.Unmarshal(file, &keys)

	return keys, unmarshalError
}
