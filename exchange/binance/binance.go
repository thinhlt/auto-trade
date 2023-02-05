package binance

import (
	"anidiot.com/auto-trade/strategy"
	"fmt"
	"sync"
	"time"

	"anidiot.com/auto-trade/entity"
	"anidiot.com/auto-trade/factory"
	"anidiot.com/auto-trade/indicator"
	"anidiot.com/auto-trade/log"
	binExchange "github.com/adshao/go-binance/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var coinMap map[string]entity.Snapshot
var processingCoin map[string]bool
var interval string
var channelID int64
var warningID int64
var mux = &sync.RWMutex{}
var resetChan chan bool

type KlineChart struct {
	candleList []*binExchange.Kline
}

func InitExchangeBot(coins []string) map[string]string {
	coinMap = make(map[string]entity.Snapshot, len(coins))
	watchMap := make(map[string]string, len(coins))
	processingCoin = make(map[string]bool, len(coins))

	interval = viper.GetString("setting.interval")
	channelID = viper.GetInt64("telegram.id")
	warningID = viper.GetInt64("telegram.warning_id")

	var wg sync.WaitGroup
	wg.Add(len(coins))
	for _, coin := range coins {
		go func(name string, wg *sync.WaitGroup) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					log.Logger.Error("Panic", zap.Any("recover", r))
				}
			}()
			var oldData []entity.BinCandle
			var err error
			for i := 1; i <= 5; i++ {
				oldData, err = getOldData(name)
				if err != nil {
					log.Logger.Error("[Binance] Get old data error", zap.String("coin", name), zap.Error(err))
					time.Sleep(5 * time.Second)
					continue
				}
				break
			}
			if len(oldData) == 0 {
				log.Logger.Error("[Binance] Get old data giveup", zap.String("coin", name))
				wg.Done()
				return
			}
			// cooking data
			opt := indicator.Calculate(oldData)
			opt.Name = name
			log.Logger.Debug("[Binance] Get old data",
				zap.String("coin", name), zap.Any("snapshot", oldData))

			factory.DB.Set30mCandle(name, oldData)
			previousTick := oldData[len(oldData)-1]
			mux.Lock()
			coinMap[name] = entity.Snapshot{
				BinCandle: previousTick,
				Opts:      opt,
			}
			watchMap[name] = "30m"
			mux.Unlock()
			wg.Done()
		}(coin, &wg)
	}
	wg.Wait()
	log.Logger.Info("Get old data done", zap.Any("coin map", watchMap))
	return watchMap
}

func KlineHandler(event *binExchange.WsKlineEvent) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	coin, ok := coinMap[event.Symbol]
	updatedBar := entity.CastWsCandle(event.Kline)
	var previous entity.Snapshot
	if !ok {
		log.Logger.Info("Not found in coin map", zap.String("coin name", event.Symbol))
		candle := entity.BinCandle{
			Name: event.Symbol,
		}
		candle.KLine = entity.KLine{
			OpenTime:         updatedBar.StartTime,
			Open:             updatedBar.Open,
			High:             updatedBar.High,
			Low:              updatedBar.Low,
			Close:            updatedBar.Close,
			Volume:           updatedBar.Volume,
			CloseTime:        updatedBar.EndTime,
			QuoteAssetVolume: updatedBar.QuoteVolume,
			TradeNum:         updatedBar.TradeNum,
		}
		mux.Lock()
		coinMap[event.Symbol] = entity.Snapshot{BinCandle: candle}
		mux.Unlock()
		return
	}
	newBar := coin.IsFinal
	if newBar {
		previous = coin
		coin = entity.Snapshot{}
		coin.Name = event.Symbol
		coin.Opts = previous.Opts
	}
	coin.OpenTime = event.Kline.StartTime
	coin.CloseTime = event.Kline.EndTime
	coin.Close = updatedBar.Close
	coin.High = updatedBar.High
	coin.Low = updatedBar.Low
	coin.Close = updatedBar.Close
	coin.Volume = updatedBar.Volume
	coin.QuoteAssetVolume = updatedBar.QuoteVolume
	coin.TradeNum = event.Kline.TradeNum
	coin.TakerBuyBaseAssetVolume = updatedBar.ActiveBuyVolume
	coin.TakerBuyQuoteAssetVolume = updatedBar.ActiveBuyQuoteVolume
	coin.CurrentTime = event.Time
	var status []string
	if newBar {
		coin.Open = updatedBar.Open
		queryResult, err := factory.DB.Get30mCandle(event.Symbol, coin.OpenTime, 28)
		if err != nil {
			log.Logger.Error("Error when get old data", zap.Error(err))
			return
		}
		oldData := append(queryResult, previous.BinCandle)
		indicator.CalculateNewTicker(event.Symbol, &coin.BinCandle, oldData, &coin.Opts)
		status = strategy.RunStrategies(event.Symbol, coin.BinCandle, &coin.Opts)
		mux.Lock()
		coinMap[event.Symbol] = coin
		mux.Unlock()
		factory.Chat.SendMessage(status, channelID)
		return
	}
	indicator.CalculateCurrent(event.Symbol, &coin.BinCandle, &coin.Opts)
	status = strategy.RunStrategies(event.Symbol, coin.BinCandle, &coin.Opts)

	if len(status) > 0 {
		mux.Lock()
		coinMap[event.Symbol] = coin
		mux.Unlock()
	}
	factory.Chat.SendMessage(status, channelID)
	if event.Kline.IsFinal {
		log.Logger.Info("final stick to store to db", zap.String("name", event.Symbol), zap.Any("kdj", coin.KDJ))
		factory.DB.SetLast30mCandle(event.Symbol, coin.BinCandle)
		coin.IsFinal = true
		mux.Lock()
		coinMap[event.Symbol] = coin
		mux.Unlock()
	}
}

func ErrHandler(err error) {
	log.Logger.Error("Binance stream", zap.Error(err))
}

func Simple(list map[string]string) {
	binExchange.WebsocketKeepalive = true
	go func() {
		for {
			doneC, _, err := binExchange.WsCombinedKlineServe(list,
				KlineHandler, ErrHandler)
			if err != nil {
				log.Logger.Error("listen livestream error", zap.Error(err))
				time.Sleep(time.Second * 3)
				continue
			}
			<-doneC
			log.Logger.Info("listen livestream done")
			time.Sleep(time.Second * 3)
		}
	}()

}
