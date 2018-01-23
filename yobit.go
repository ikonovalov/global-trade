package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"strconv"
	"net/url"
	"bytes"
	"github.com/ikonovalov/go-cloudflare-scraper"
)

const (
	ApiBase    = "https://yobit.net/api/"
	ApiVersion = "3"
	ApiTrade   = "https://yobit.net/tapi/"
)

func NewYobit() (*Yobit) {
	cloudflare, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		panic(err)
	}

	keys, err := loadApiKeys()
	if err != nil {
		panic(err)
	}

	yobit := Yobit{client: &http.Client{Transport: cloudflare, Jar: cloudflare.Cookies}, apiKeys: &keys}
	yobit.PassCloudflare()
	return &yobit
}

func (y *Yobit) PassCloudflare() (*Yobit) {
	channel := make(chan TickerInfoResponse)
	go y.Tickers24("btc_usd", channel)
	<-channel
	return y
}

func (y *Yobit) Tickers24(pairs string, ch chan TickerInfoResponse) {
	ticker24Url := ApiBase + ApiVersion + "/ticker/" + pairs
	response := y.callPublic(ticker24Url)

	var tickerResponse TickerInfoResponse
	pTicker := &tickerResponse.Tickers

	if err := unmarshal(response, pTicker); err != nil {
		panic(err)
	}
	ch <- tickerResponse
}

func (y *Yobit) Info(ch chan InfoResponse) {
	infoUrl := ApiBase + ApiVersion + "/info"
	response := y.callPublic(infoUrl)
	var infoResponse InfoResponse
	if err := unmarshal(response, &infoResponse); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- infoResponse
}

func (y *Yobit) Depth(pairs string, ch chan DepthResponse) {
	y.DepthLimited(pairs, 150, ch)
}

func (y *Yobit) DepthLimited(pairs string, limit int, ch chan DepthResponse) {
	limitedDepthUrl := fmt.Sprintf("%s/depth/%s?limit=%d", ApiBase+ApiVersion, pairs, limit)
	response := y.callPublic(limitedDepthUrl)
	var depthResponse DepthResponse
	if err := unmarshal(response, &depthResponse.Orders); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- depthResponse
}

func (y *Yobit) TradesLimited(pairs string, limit int, ch chan TradesResponse) {
	tradesLimitedUrl := fmt.Sprintf("%s/trades/%s?limit=%d", ApiBase+ApiVersion, pairs, limit)
	response := y.callPublic(tradesLimitedUrl)
	var tradesResponse TradesResponse
	if err := unmarshal(response, &tradesResponse.Trades); err != nil {
		log.Fatal(err)
		panic(err)
	}
	ch <- tradesResponse
}

func (y *Yobit) GetInfo(ch chan GetInfoResponse) {
	response := y.callPrivate("getInfo")
	var getInfoResp GetInfoResponse
	if err := unmarshal(response, &getInfoResp); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if getInfoResp.Success == 0 {
		panic(errors.New(getInfoResp.Error))
	}
	ch <- getInfoResp
}

func unmarshal(data [] byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unmarshaling failed\n%s\n%s", string(data), err.Error()))
	}
	return err
}

func (y *Yobit) query(req *http.Request) ([]byte) {
	resp, err := y.client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("error. HTTP %d", resp.StatusCode)
		panic(errors.New(fmt.Sprintf("\n%s\nSomething goes wrong. We got HTTP%d", req.URL.String(), resp.StatusCode)))
	}
	response, _ := ioutil.ReadAll(resp.Body)
	return response
}

func (y *Yobit) callPublic(url string) ([]byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic(err)
	}
	return y.query(req)
}

func (y *Yobit) callPrivate(method string) ([]byte) {
	nonce := y.GetNonce()
	form := url.Values{
		"method": {method},
		"nonce":  {strconv.FormatUint(nonce, 10)},
	}
	encode := form.Encode()
	signature := signHmacSha512([]byte(y.apiKeys.Secret), []byte(encode))
	body := bytes.NewBufferString(encode)
	req, err := http.NewRequest("POST", ApiTrade, body)

	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic(err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", y.apiKeys.Key)
	req.Header.Add("Sign", signature)

	query := y.query(req)
	incrementNonce(&nonce)
	return query
}

type Yobit struct {
	client  *http.Client
	apiKeys *ApiKeys
}
