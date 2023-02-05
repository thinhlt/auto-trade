package factory

import (
	"anidiot.com/auto-trade/bot"
	msqlClient "anidiot.com/auto-trade/db/msql-client"
	"anidiot.com/auto-trade/entity"
	"anidiot.com/auto-trade/log"
	"anidiot.com/common/kafka"
	"github.com/spf13/viper"
)

var Chat ChatBot
var DB StoreDB
var Producer kafka.Producer
var WarningMap map[string]struct{}

type ChatBot interface {
	SendMessage(messages []string, channelID int64)
	Close()
}

type StoreDB interface {
	Set1hCandle(candle entity.CoinTicker)
	Set1hCandleWithTime(candle entity.CoinTicker, timeString string)
	SetLast1hCandle(name string, candle entity.Candle, score float64)
	SetCandle(key string, candles []entity.Candle, start int64)
	Set15mCandle(name string, candles []entity.ICandle)
	Set30mCandle(name string, candles []entity.BinCandle)
	SetLast15mCandle(name string, candle entity.ICandle)
	SetLast30mCandle(name string, candle entity.BinCandle)
	Get15mCandle(name string, currentTime int64, num int) ([]entity.ICandle, error)
	Get30mCandle(name string, currentTime int64, num int) ([]entity.BinCandle, error)
	Destroy()
}

func InitProduction() {
	WarningMap = make(map[string]struct{}, 0)
	warningList := viper.GetStringSlice("setting.warning_list")
	for _, pair := range warningList {
		WarningMap[pair] = struct{}{}
	}
	chatbot := bot.InitBot(WarningMap)
	Chat = chatbot
	DB = msqlClient.Init()
	Producer = kafka.SimpleProducer(viper.GetStringSlice("kafka.broker"),
		viper.GetString("kafka.topic"))
}

func Destroy() {
	DB.Destroy()
	Chat.Close()
}

func InitDevelopment() {
	Chat = log.LogStore{}
	DB = log.LogStore{}
}
