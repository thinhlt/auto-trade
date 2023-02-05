package bot

import (
	"anidiot.com/auto-trade/utils"
)

type CoinPair struct {
	Name          string
	Percent       float64
	BuyPrice      float64
	SellPrice     float64
	UpperStoploss float64
	LowerStoploss float64
	MiddleStop    float64
	BaseVolume    float64
	QuoteVolume   float64
	CVD           []int64
	CVDQueue      utils.QueueInt64
	Have          bool
	Level         int
}

type Ticker struct {
	Name              string  `json:"MarketName,omitempty"`
	High              float64 `json:"High,omitempty"`
	Low               float64 `json:"Low,omitempty"`
	Volume            float64 `json:"Volume,omitempty"`
	Last              float64 `json:"Last,omitempty"`
	BaseVolume        float64 `json:"BaseVolume,omitempty"`
	TimeStamp         string  `json:"TimeStamp,omitempty"`
	Bid               float64 `json:"Bid,omitempty"`
	Ask               float64 `json:"Ask,omitempty"`
	OpenBuyOrders     int16   `json:"OpenBuyOrders,omitempty"`
	OpenSellOrders    int16   `json:"OpenSellOrders,omitempty"`
	PrevDay           float64 `json:"PrevDay,omitempty"`
	Created           string  `json:"Created,omitempty"`
	DisplayMarketName string  `json:"DisplayMarketName,omitempty"`
}

type MarketSummary struct {
	Status  bool     `json:"success,omitempty"`
	Message string   `json:"message,omitempty"`
	Result  []Ticker `json:"result,omitempty"`
}
