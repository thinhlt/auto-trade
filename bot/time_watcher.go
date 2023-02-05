package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"log"
)

var timer <-chan time.Time
var closeSignal <-chan bool
var tickerChan chan Ticker
var client http.Client
var watchingCoins map[string]CoinPair

func (bot Bot) initWatcher(duration time.Duration) {
	ticker := time.NewTicker(duration)
	client = http.Client{
		Timeout: time.Second * 2,
	}
	tickerChan = make(chan Ticker, 10)
	bot.wg.Add(1)
	for k := range watchingCoins {
		go checkMarket(k)
	}
	go func() {
		fmt.Println("Init time watcher")
		defer bot.wg.Done()
		for {
			select {
			case signal := <-bot.stopChan:
				if signal == "timer" {
					fmt.Printf("stop timer")
					ticker.Stop()
					return
				} else {
					bot.stopChan <- signal
				}
			case t := <-ticker.C:
				fmt.Println("ticker fire", t)
				for k := range watchingCoins {
					go checkMarket(k)
				}
				if coinBit, err := json.Marshal(watchingCoins); err != nil {
					fmt.Println("coin error", err)
				} else {
					fmt.Println("coin", string(coinBit))
				}
			}
		}
	}()
}

func checkMarket(market string) {
	// log.Println("check market", market)
	url := fmt.Sprintf("https://api.entity.com/api/v1.1/public/getmarketsummary?market=%s", market)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("request market error", url, err)
		return
	}
	res, resErr := client.Do(request)
	if resErr != nil {
		log.Printf("response error", url, err)
		return
	}
	if res == nil {
		log.Printf("response empty", url)
		return
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Printf("cannot read body", url)
		return
	}
	summary := MarketSummary{}
	jsonErr := json.Unmarshal(body, &summary)
	if jsonErr != nil {
		log.Printf("parse request fail", url, jsonErr)
		return
	}
	tickerChan <- summary.Result[0]
}
