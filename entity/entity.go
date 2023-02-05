package entity

type ICandle interface {
	ToByteArray() []byte
	GetTime() int64
	GetScore() uint
}

type Trend int

const (
	NOTREND Trend = iota
	UPTREND
	DOWNTREND
	STOP
)

func (t Trend) ToString() string {
	return trendMap[t]
}

var trendMap = map[Trend]string{
	UPTREND:   "UPTREND",
	DOWNTREND: "DOWNTREND",
	STOP:      "STOP",
	NOTREND:   "NOTREND",
}

type Market struct {
	Name   string `json:"symbol,omitempty"`
	Status string `json:"status,omitempty"`
}

type CoinTicker struct {
	Name     string
	M15Chart []Candle
	H1Chart  []Candle
}

type Candle struct {
	Hight         float64       `json:"high,string,omitempty"`
	Low           float64       `json:"low,string,omitempty"`
	Open          float64       `json:"open,string,omitempty"`
	Last          float64       `json:"last,string,omitempty"`
	Close         float64       `json:"close,string,omitempty"`
	StartAt       CustomTime    `json:"startsAt,omitempty"`
	Volume        float64       `json:"volume,string,omitempty"`
	QuoteVolume   float64       `json:"quoteVolume,string,omitempty"`
	EMA           float64       `json:"ema,omitempty"`
	Psar          float64       `json:"psar,omitempty"`
	StochasticRSI float64       `json:"stochastic,omitempty"`
	BBB           BoilingerBand `json:"bbb,omitempty"`
	RSI           float64       `json:"rsi,omitempty"`
}

type BoilingerBand struct {
	Basic float64
	Upper float64
	Lower float64
}

type PrevEMA struct {
	EMA12 float64
	EMA25 float64
	EMA50 float64
}

type Ticker struct {
	MarketName string      `json:"MarketName,omitempty"`
	Name       string      `json:"Name,omitempty"`
	Hight      float64     `json:"High,omitempty"`
	Low        float64     `json:"Low,omitempty"`
	Last       float64     `json:"Last,omitempty"`
	StartAt    SummaryTime `json:"TimeStamp,omitempty"`
}

type MarketSummary struct {
	Status  bool     `json:"success,omitempty"`
	Message string   `json:"message,omitempty"`
	Result  []Ticker `json:"result,omitempty"`
}

type HaCounter struct {
	Trend
	Count int
}

type ChartOption struct {
	High               float64   `json:"high,omitempty"`
	Low                float64   `json:"low,omitempty"`
	AF                 float64   `json:"af,omitempty"`
	BBBState           BBState   `json:"state,omitempty"`
	TrendStopLost      Trend
	HaStopLost         float64
	PrevKDJ            KDJ
	Name               string
	Trend              Trend
	High9              float64
	Low9               float64
	PrevEMAList        PrevEMA
}

type KDJ struct {
	K float64 `json:"k_percent,omitempty" gorm:"-"`
	D float64 `json:"d_percent,omitempty" gorm:"-"`
	J float64 `json:"j_percent,omitempty" gorm:"-"`
}