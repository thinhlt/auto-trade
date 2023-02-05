package enitty

import "time"

type Snapshot struct {
	BinCandle
	Opts          ChartOption
	LogType ActionLog
}

type KLine struct {
	OpenTime                 int64   `json:"openTime"`
	Open                     float64 `json:"open,string"`
	High                     float64 `json:"high,string"`
	Low                      float64 `json:"low,string"`
	Close                    float64 `json:"close,string"`
	CloseTime                int64   `json:"closeTime"`
	TradeNum                 int64   `json:"tradeNum,string"`
}

type HA struct {
	HAOpen  float64 `json:"HAopen,string"`
	HAHigh  float64 `json:"HAhigh,string"`
	HALow   float64 `json:"HAlow,string"`
	HAClose float64 `json:"HAclose,string"`
}

type ADX struct {
	ADX   float64 `json:"adx,string"`
	DIP14 float64 `json:"dip14,string"`
	DIM14 float64 `json:"dim14,string"`
	DIP   float64 `json:"dip,string"`
	DIM   float64 `json:"dim,string"`
	TR    float64 `json:"tr,string"`
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
	EMA26             float64       `json:"ema26,omitempty" gorm:"-"`
	EMA12             float64       `json:"ema12,omitempty" gorm:"-"`
	EMA200            float64       `json:"ema200,omitempty" gorm:"-"`
	Psar              float64       `json:"psar,omitempty"`
	BBB               BoilingerBand `json:"bbb,omitempty" gorm:"-"`
	RSI               float64       `json:"rsi,omitempty" gorm:"-"`
	Score             uint          `json:"score,omitempty" gorm:"primaryKey"`
	CurrentTime       int64         `gorm:"-"`
	Name              string        `json:"name,omitempty" gorm:"-"`
}

type ChartOption struct {
	High             float64   `json:"high,omitempty"`
	Low              float64   `json:"low,omitempty"`
	AF               float64   `json:"af,omitempty"`
	BBBState         BBState   `json:"state,omitempty"`
	StopAt           time.Time `json:"startsAt,omitempty"`
	Name             string
	Trend            Trend
	PreviousPsarDiff float64
	PredictPsar      float64
	PsarLevel        float64
}

type BoilingerBand struct {
	Basic float64
	Upper float64
	Lower float64
}

type BBState int

const (
	None BBState = iota
	UnderLower
	UnderBasic
	UnderUpper
	OverUpper
)

var bbStateMap = map[BBState]string{
	None:       "",
	UnderLower: "underLower",
	UnderBasic: "underBasic",
	UnderUpper: "underUpper",
	OverUpper:  "overUpper",
}

func (state BBState) ToString() string {
	return bbStateMap[state]
}


type Trend int

const (
	UPTREND   Trend = iota
	DOWNTREND Trend = iota
)

var trendMap = map[Trend]string{
	UPTREND:   "UPTREND",
	DOWNTREND: "DOWNTREND",
}

func (t Trend) ToString() string {
	return trendMap[t]
}

type ActionLog int

const (
	PSARCROSS   ActionLog = iota
	PSARSTRONG ActionLog = iota
	BBMID ActionLog = iota
	BBPOLE ActionLog = iota
)

var actionMap = map[ActionLog]string{
	PSARCROSS:   "PSARCROSS",
	PSARSTRONG: "PSARSTRONG",
	BBMID: "BBMID",
	BBPOLE: "BBPOLE",
}

func (t ActionLog) ToString() string {
	return actionMap[t]
}
