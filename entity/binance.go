package entity

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2"
)

var TIME_PARTITION int64 = 1800000

type Snapshot struct {
	BinCandle
	Opts          ChartOption
	IsFinal       bool
}

type KLine struct {
	OpenTime                 int64   `json:"openTime"`
	Open                     float64 `json:"open,string"`
	High                     float64 `json:"high,string"`
	Low                      float64 `json:"low,string"`
	Close                    float64 `json:"close,string"`
	Volume                   float64 `json:"volume,string"`
	CloseTime                int64   `json:"closeTime"`
	QuoteAssetVolume         float64 `json:"quoteAssetVolume,string"`
	TradeNum                 int64   `json:"tradeNum,string"`
	TakerBuyBaseAssetVolume  float64 `json:"takerBuyBaseAssetVolume,string"`
	TakerBuyQuoteAssetVolume float64 `json:"takerBuyQuoteAssetVolume,string"`
}

type HA struct {
	HAOpen  float64 `json:"HAopen,string"`
	HAHigh  float64 `json:"HAhigh,string"`
	HALow   float64 `json:"HAlow,string"`
	HAClose float64 `json:"HAclose,string"`
}

// WsKline define websocket kline
type WsKline struct {
	StartTime            int64   `json:"t"`
	EndTime              int64   `json:"T"`
	Symbol               string  `json:"s"`
	Interval             string  `json:"i"`
	Open                 float64 `json:"o,string"`
	Close                float64 `json:"c,string"`
	High                 float64 `json:"h,string"`
	Low                  float64 `json:"l,string"`
	Volume               float64 `json:"v,string"`
	TradeNum             int64   `json:"n"`
	IsFinal              bool    `json:"x"`
	QuoteVolume          float64 `json:"q,string"`
	ActiveBuyVolume      float64 `json:"V,string"`
	ActiveBuyQuoteVolume float64 `json:"Q,string"`
}

type BinCandle struct {
	KLine
	HA
	Psar              float64       `json:"psar,omitempty"`
	BBB               BoilingerBand `json:"bbb,omitempty" gorm:"-"`
	RSI               float64       `json:"rsi,omitempty" gorm:"-"`
	VolumeSMA21       float64       `json:"VolumeSMA21,omitempty" gorm:"-"`
	Score             uint          `json:"score,omitempty" gorm:"primaryKey"`
	CurrentTime       int64         `gorm:"-"`
	Name              string        `json:"name,omitempty" gorm:"-"`
	KDJ
	State   OrderState `json:"state,omitempty" gorm:"-"`
	EMA12   float64    `gorm:"-"`
	EMA25   float64    `gorm:"-"`
	EMA50   float64    `gorm:"-"`
}

func CastCandle(kline *binance.Kline) BinCandle {
	var result BinCandle
	jsonData, _ := json.Marshal(kline)
	json.Unmarshal(jsonData, &result)
	return result
}

func CastWsCandle(kline binance.WsKline) WsKline {
	var result WsKline
	jsonData, _ := json.Marshal(kline)
	json.Unmarshal(jsonData, &result)
	return result
}

func (b BinCandle) ToByteArray() []byte {
	result, err := json.Marshal(b)
	if err != nil {
		//log.Logger.Error("error when marshall candle", zap.Any("candle", b), zap.Error(err))
		return nil
	}
	return result
}

func (b BinCandle) GetTime() int64 {
	return b.OpenTime
}

func (b BinCandle) GetScore() uint {
	if b.Score == 0 {
		b.Score = uint(b.OpenTime / TIME_PARTITION)
	}
	return b.Score
}

func (b BinCandle) GetDayString() string {
	return time.Unix(b.OpenTime, 0).Format("0601")
}
